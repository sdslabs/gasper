package shell

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/operation"
	"github.com/chrislusf/seaweedfs/weed/pb/master_pb"
	"github.com/chrislusf/seaweedfs/weed/pb/volume_server_pb"
	"github.com/chrislusf/seaweedfs/weed/storage/erasure_coding"
	"github.com/chrislusf/seaweedfs/weed/storage/needle"
	"google.golang.org/grpc"
)

func moveMountedShardToEcNode(commandEnv *CommandEnv, existingLocation *EcNode, collection string, vid needle.VolumeId, shardId erasure_coding.ShardId, destinationEcNode *EcNode, applyBalancing bool) (err error) {

	copiedShardIds := []uint32{uint32(shardId)}

	if applyBalancing {

		// ask destination node to copy shard and the ecx file from source node, and mount it
		copiedShardIds, err = oneServerCopyAndMountEcShardsFromSource(commandEnv.option.GrpcDialOption, destinationEcNode, []uint32{uint32(shardId)}, vid, collection, existingLocation.info.Id)
		if err != nil {
			return err
		}

		// unmount the to be deleted shards
		err = unmountEcShards(commandEnv.option.GrpcDialOption, vid, existingLocation.info.Id, copiedShardIds)
		if err != nil {
			return err
		}

		// ask source node to delete the shard, and maybe the ecx file
		err = sourceServerDeleteEcShards(commandEnv.option.GrpcDialOption, collection, vid, existingLocation.info.Id, copiedShardIds)
		if err != nil {
			return err
		}

		fmt.Printf("moved ec shard %d.%d %s => %s\n", vid, shardId, existingLocation.info.Id, destinationEcNode.info.Id)

	}

	destinationEcNode.addEcVolumeShards(vid, collection, copiedShardIds)
	existingLocation.deleteEcVolumeShards(vid, copiedShardIds)

	return nil

}

func oneServerCopyAndMountEcShardsFromSource(grpcDialOption grpc.DialOption,
	targetServer *EcNode, shardIdsToCopy []uint32,
	volumeId needle.VolumeId, collection string, existingLocation string) (copiedShardIds []uint32, err error) {

	fmt.Printf("allocate %d.%v %s => %s\n", volumeId, shardIdsToCopy, existingLocation, targetServer.info.Id)

	err = operation.WithVolumeServerClient(targetServer.info.Id, grpcDialOption, func(volumeServerClient volume_server_pb.VolumeServerClient) error {

		if targetServer.info.Id != existingLocation {

			fmt.Printf("copy %d.%v %s => %s\n", volumeId, shardIdsToCopy, existingLocation, targetServer.info.Id)
			_, copyErr := volumeServerClient.VolumeEcShardsCopy(context.Background(), &volume_server_pb.VolumeEcShardsCopyRequest{
				VolumeId:       uint32(volumeId),
				Collection:     collection,
				ShardIds:       shardIdsToCopy,
				CopyEcxFile:    true,
				CopyEcjFile:    true,
				CopyVifFile:    true,
				SourceDataNode: existingLocation,
			})
			if copyErr != nil {
				return fmt.Errorf("copy %d.%v %s => %s : %v\n", volumeId, shardIdsToCopy, existingLocation, targetServer.info.Id, copyErr)
			}
		}

		fmt.Printf("mount %d.%v on %s\n", volumeId, shardIdsToCopy, targetServer.info.Id)
		_, mountErr := volumeServerClient.VolumeEcShardsMount(context.Background(), &volume_server_pb.VolumeEcShardsMountRequest{
			VolumeId:   uint32(volumeId),
			Collection: collection,
			ShardIds:   shardIdsToCopy,
		})
		if mountErr != nil {
			return fmt.Errorf("mount %d.%v on %s : %v\n", volumeId, shardIdsToCopy, targetServer.info.Id, mountErr)
		}

		if targetServer.info.Id != existingLocation {
			copiedShardIds = shardIdsToCopy
			glog.V(0).Infof("%s ec volume %d deletes shards %+v", existingLocation, volumeId, copiedShardIds)
		}

		return nil
	})

	if err != nil {
		return
	}

	return
}

func eachDataNode(topo *master_pb.TopologyInfo, fn func(dc string, rack RackId, dn *master_pb.DataNodeInfo)) {
	for _, dc := range topo.DataCenterInfos {
		for _, rack := range dc.RackInfos {
			for _, dn := range rack.DataNodeInfos {
				fn(dc.Id, RackId(rack.Id), dn)
			}
		}
	}
}

func sortEcNodesByFreeslotsDecending(ecNodes []*EcNode) {
	sort.Slice(ecNodes, func(i, j int) bool {
		return ecNodes[i].freeEcSlot > ecNodes[j].freeEcSlot
	})
}

func sortEcNodesByFreeslotsAscending(ecNodes []*EcNode) {
	sort.Slice(ecNodes, func(i, j int) bool {
		return ecNodes[i].freeEcSlot < ecNodes[j].freeEcSlot
	})
}

type CandidateEcNode struct {
	ecNode     *EcNode
	shardCount int
}

// if the index node changed the freeEcSlot, need to keep every EcNode still sorted
func ensureSortedEcNodes(data []*CandidateEcNode, index int, lessThan func(i, j int) bool) {
	for i := index - 1; i >= 0; i-- {
		if lessThan(i+1, i) {
			swap(data, i, i+1)
		} else {
			break
		}
	}
	for i := index + 1; i < len(data); i++ {
		if lessThan(i, i-1) {
			swap(data, i, i-1)
		} else {
			break
		}
	}
}

func swap(data []*CandidateEcNode, i, j int) {
	t := data[i]
	data[i] = data[j]
	data[j] = t
}

func countShards(ecShardInfos []*master_pb.VolumeEcShardInformationMessage) (count int) {
	for _, ecShardInfo := range ecShardInfos {
		shardBits := erasure_coding.ShardBits(ecShardInfo.EcIndexBits)
		count += shardBits.ShardIdCount()
	}
	return
}

func countFreeShardSlots(dn *master_pb.DataNodeInfo) (count int) {
	return int(dn.MaxVolumeCount-dn.ActiveVolumeCount)*erasure_coding.DataShardsCount - countShards(dn.EcShardInfos)
}

type RackId string
type EcNodeId string

type EcNode struct {
	info       *master_pb.DataNodeInfo
	dc         string
	rack       RackId
	freeEcSlot int
}

func (ecNode *EcNode) localShardIdCount(vid uint32) int {
	for _, ecShardInfo := range ecNode.info.EcShardInfos {
		if vid == ecShardInfo.Id {
			shardBits := erasure_coding.ShardBits(ecShardInfo.EcIndexBits)
			return shardBits.ShardIdCount()
		}
	}
	return 0
}

type EcRack struct {
	ecNodes    map[EcNodeId]*EcNode
	freeEcSlot int
}

func collectEcNodes(commandEnv *CommandEnv, selectedDataCenter string) (ecNodes []*EcNode, totalFreeEcSlots int, err error) {

	// list all possible locations
	var resp *master_pb.VolumeListResponse
	err = commandEnv.MasterClient.WithClient(func(client master_pb.SeaweedClient) error {
		resp, err = client.VolumeList(context.Background(), &master_pb.VolumeListRequest{})
		return err
	})
	if err != nil {
		return nil, 0, err
	}

	// find out all volume servers with one slot left.
	ecNodes, totalFreeEcSlots = collectEcVolumeServersByDc(resp.TopologyInfo, selectedDataCenter)

	sortEcNodesByFreeslotsDecending(ecNodes)

	return
}

func collectEcVolumeServersByDc(topo *master_pb.TopologyInfo, selectedDataCenter string) (ecNodes []*EcNode, totalFreeEcSlots int) {
	eachDataNode(topo, func(dc string, rack RackId, dn *master_pb.DataNodeInfo) {
		if selectedDataCenter != "" && selectedDataCenter != dc {
			return
		}

		freeEcSlots := countFreeShardSlots(dn)
		ecNodes = append(ecNodes, &EcNode{
			info:       dn,
			dc:         dc,
			rack:       rack,
			freeEcSlot: int(freeEcSlots),
		})
		totalFreeEcSlots += freeEcSlots
	})
	return
}

func sourceServerDeleteEcShards(grpcDialOption grpc.DialOption, collection string, volumeId needle.VolumeId, sourceLocation string, toBeDeletedShardIds []uint32) error {

	fmt.Printf("delete %d.%v from %s\n", volumeId, toBeDeletedShardIds, sourceLocation)

	return operation.WithVolumeServerClient(sourceLocation, grpcDialOption, func(volumeServerClient volume_server_pb.VolumeServerClient) error {
		_, deleteErr := volumeServerClient.VolumeEcShardsDelete(context.Background(), &volume_server_pb.VolumeEcShardsDeleteRequest{
			VolumeId:   uint32(volumeId),
			Collection: collection,
			ShardIds:   toBeDeletedShardIds,
		})
		return deleteErr
	})

}

func unmountEcShards(grpcDialOption grpc.DialOption, volumeId needle.VolumeId, sourceLocation string, toBeUnmountedhardIds []uint32) error {

	fmt.Printf("unmount %d.%v from %s\n", volumeId, toBeUnmountedhardIds, sourceLocation)

	return operation.WithVolumeServerClient(sourceLocation, grpcDialOption, func(volumeServerClient volume_server_pb.VolumeServerClient) error {
		_, deleteErr := volumeServerClient.VolumeEcShardsUnmount(context.Background(), &volume_server_pb.VolumeEcShardsUnmountRequest{
			VolumeId: uint32(volumeId),
			ShardIds: toBeUnmountedhardIds,
		})
		return deleteErr
	})
}

func mountEcShards(grpcDialOption grpc.DialOption, collection string, volumeId needle.VolumeId, sourceLocation string, toBeMountedhardIds []uint32) error {

	fmt.Printf("mount %d.%v on %s\n", volumeId, toBeMountedhardIds, sourceLocation)

	return operation.WithVolumeServerClient(sourceLocation, grpcDialOption, func(volumeServerClient volume_server_pb.VolumeServerClient) error {
		_, mountErr := volumeServerClient.VolumeEcShardsMount(context.Background(), &volume_server_pb.VolumeEcShardsMountRequest{
			VolumeId:   uint32(volumeId),
			Collection: collection,
			ShardIds:   toBeMountedhardIds,
		})
		return mountErr
	})
}

func divide(total, n int) float64 {
	return float64(total) / float64(n)
}

func ceilDivide(total, n int) int {
	return int(math.Ceil(float64(total) / float64(n)))
}

func findEcVolumeShards(ecNode *EcNode, vid needle.VolumeId) erasure_coding.ShardBits {

	for _, shardInfo := range ecNode.info.EcShardInfos {
		if needle.VolumeId(shardInfo.Id) == vid {
			return erasure_coding.ShardBits(shardInfo.EcIndexBits)
		}
	}

	return 0
}

func (ecNode *EcNode) addEcVolumeShards(vid needle.VolumeId, collection string, shardIds []uint32) *EcNode {

	foundVolume := false
	for _, shardInfo := range ecNode.info.EcShardInfos {
		if needle.VolumeId(shardInfo.Id) == vid {
			oldShardBits := erasure_coding.ShardBits(shardInfo.EcIndexBits)
			newShardBits := oldShardBits
			for _, shardId := range shardIds {
				newShardBits = newShardBits.AddShardId(erasure_coding.ShardId(shardId))
			}
			shardInfo.EcIndexBits = uint32(newShardBits)
			ecNode.freeEcSlot -= newShardBits.ShardIdCount() - oldShardBits.ShardIdCount()
			foundVolume = true
			break
		}
	}

	if !foundVolume {
		var newShardBits erasure_coding.ShardBits
		for _, shardId := range shardIds {
			newShardBits = newShardBits.AddShardId(erasure_coding.ShardId(shardId))
		}
		ecNode.info.EcShardInfos = append(ecNode.info.EcShardInfos, &master_pb.VolumeEcShardInformationMessage{
			Id:          uint32(vid),
			Collection:  collection,
			EcIndexBits: uint32(newShardBits),
		})
		ecNode.freeEcSlot -= len(shardIds)
	}

	return ecNode
}

func (ecNode *EcNode) deleteEcVolumeShards(vid needle.VolumeId, shardIds []uint32) *EcNode {

	for _, shardInfo := range ecNode.info.EcShardInfos {
		if needle.VolumeId(shardInfo.Id) == vid {
			oldShardBits := erasure_coding.ShardBits(shardInfo.EcIndexBits)
			newShardBits := oldShardBits
			for _, shardId := range shardIds {
				newShardBits = newShardBits.RemoveShardId(erasure_coding.ShardId(shardId))
			}
			shardInfo.EcIndexBits = uint32(newShardBits)
			ecNode.freeEcSlot -= newShardBits.ShardIdCount() - oldShardBits.ShardIdCount()
		}
	}

	return ecNode
}

func groupByCount(data []*EcNode, identifierFn func(*EcNode) (id string, count int)) map[string]int {
	countMap := make(map[string]int)
	for _, d := range data {
		id, count := identifierFn(d)
		countMap[id] += count
	}
	return countMap
}

func groupBy(data []*EcNode, identifierFn func(*EcNode) (id string)) map[string][]*EcNode {
	groupMap := make(map[string][]*EcNode)
	for _, d := range data {
		id := identifierFn(d)
		groupMap[id] = append(groupMap[id], d)
	}
	return groupMap
}
