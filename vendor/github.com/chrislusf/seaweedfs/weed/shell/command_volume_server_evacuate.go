package shell

import (
	"context"
	"flag"
	"fmt"
	"github.com/chrislusf/seaweedfs/weed/pb/master_pb"
	"github.com/chrislusf/seaweedfs/weed/storage/erasure_coding"
	"github.com/chrislusf/seaweedfs/weed/storage/needle"
	"github.com/chrislusf/seaweedfs/weed/storage/super_block"
	"io"
	"sort"
)

func init() {
	Commands = append(Commands, &commandVolumeServerEvacuate{})
}

type commandVolumeServerEvacuate struct {
}

func (c *commandVolumeServerEvacuate) Name() string {
	return "volumeServer.evacuate"
}

func (c *commandVolumeServerEvacuate) Help() string {
	return `move out all data on a volume server

	volumeServer.evacuate -node <host:port>

	This command moves all data away from the volume server.
	The volumes on the volume servers will be redistributed.

	Usually this is used to prepare to shutdown or upgrade the volume server.

	Sometimes a volume can not be moved because there are no
	good destination to meet the replication requirement. 
	E.g. a volume replication 001 in a cluster with 2 volume servers can not be moved.
	You can use "-skipNonMoveable" to move the rest volumes.

`
}

func (c *commandVolumeServerEvacuate) Do(args []string, commandEnv *CommandEnv, writer io.Writer) (err error) {

	if err = commandEnv.confirmIsLocked(); err != nil {
		return
	}

	vsEvacuateCommand := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	volumeServer := vsEvacuateCommand.String("node", "", "<host>:<port> of the volume server")
	skipNonMoveable := vsEvacuateCommand.Bool("skipNonMoveable", false, "skip volumes that can not be moved")
	applyChange := vsEvacuateCommand.Bool("force", false, "actually apply the changes")
	if err = vsEvacuateCommand.Parse(args); err != nil {
		return nil
	}

	if *volumeServer == "" {
		return fmt.Errorf("need to specify volume server by -node=<host>:<port>")
	}

	return volumeServerEvacuate(commandEnv, *volumeServer, *skipNonMoveable, *applyChange, writer)

}

func volumeServerEvacuate(commandEnv *CommandEnv, volumeServer string, skipNonMoveable, applyChange bool, writer io.Writer) (err error) {
	// 1. confirm the volume server is part of the cluster
	// 2. collect all other volume servers, sort by empty slots
	// 3. move to any other volume server as long as it satisfy the replication requirements

	// list all the volumes
	var resp *master_pb.VolumeListResponse
	err = commandEnv.MasterClient.WithClient(func(client master_pb.SeaweedClient) error {
		resp, err = client.VolumeList(context.Background(), &master_pb.VolumeListRequest{})
		return err
	})
	if err != nil {
		return err
	}

	if err := evacuateNormalVolumes(commandEnv, resp, volumeServer, skipNonMoveable, applyChange, writer); err != nil {
		return err
	}

	if err := evacuateEcVolumes(commandEnv, resp, volumeServer, skipNonMoveable, applyChange, writer); err != nil {
		return err
	}

	return nil
}

func evacuateNormalVolumes(commandEnv *CommandEnv, resp *master_pb.VolumeListResponse, volumeServer string, skipNonMoveable, applyChange bool, writer io.Writer) error {
	// find this volume server
	volumeServers := collectVolumeServersByDc(resp.TopologyInfo, "")
	thisNode, otherNodes := nodesOtherThan(volumeServers, volumeServer)
	if thisNode == nil {
		return fmt.Errorf("%s is not found in this cluster", volumeServer)
	}

	// move away normal volumes
	volumeReplicas, _ := collectVolumeReplicaLocations(resp)
	for _, vol := range thisNode.info.VolumeInfos {
		hasMoved, err := moveAwayOneNormalVolume(commandEnv, volumeReplicas, vol, thisNode, otherNodes, applyChange)
		if err != nil {
			return fmt.Errorf("move away volume %d from %s: %v", vol.Id, volumeServer, err)
		}
		if !hasMoved {
			if skipNonMoveable {
				replicaPlacement, _ := super_block.NewReplicaPlacementFromByte(byte(vol.ReplicaPlacement))
				fmt.Fprintf(writer, "skipping non moveable volume %d replication:%s\n", vol.Id, replicaPlacement.String())
			} else {
				return fmt.Errorf("failed to move volume %d from %s", vol.Id, volumeServer)
			}
		}
	}
	return nil
}

func evacuateEcVolumes(commandEnv *CommandEnv, resp *master_pb.VolumeListResponse, volumeServer string, skipNonMoveable, applyChange bool, writer io.Writer) error {
	// find this ec volume server
	ecNodes, _ := collectEcVolumeServersByDc(resp.TopologyInfo, "")
	thisNode, otherNodes := ecNodesOtherThan(ecNodes, volumeServer)
	if thisNode == nil {
		return fmt.Errorf("%s is not found in this cluster\n", volumeServer)
	}

	// move away ec volumes
	for _, ecShardInfo := range thisNode.info.EcShardInfos {
		hasMoved, err := moveAwayOneEcVolume(commandEnv, ecShardInfo, thisNode, otherNodes, applyChange)
		if err != nil {
			return fmt.Errorf("move away volume %d from %s: %v", ecShardInfo.Id, volumeServer, err)
		}
		if !hasMoved {
			if skipNonMoveable {
				fmt.Fprintf(writer, "failed to move away ec volume %d from %s\n", ecShardInfo.Id, volumeServer)
			} else {
				return fmt.Errorf("failed to move away ec volume %d from %s", ecShardInfo.Id, volumeServer)
			}
		}
	}
	return nil
}

func moveAwayOneEcVolume(commandEnv *CommandEnv, ecShardInfo *master_pb.VolumeEcShardInformationMessage, thisNode *EcNode, otherNodes []*EcNode, applyChange bool) (hasMoved bool, err error) {

	for _, shardId := range erasure_coding.ShardBits(ecShardInfo.EcIndexBits).ShardIds() {

		sort.Slice(otherNodes, func(i, j int) bool {
			return otherNodes[i].localShardIdCount(ecShardInfo.Id) < otherNodes[j].localShardIdCount(ecShardInfo.Id)
		})

		for i := 0; i < len(otherNodes); i++ {
			emptyNode := otherNodes[i]
			err = moveMountedShardToEcNode(commandEnv, thisNode, ecShardInfo.Collection, needle.VolumeId(ecShardInfo.Id), shardId, emptyNode, applyChange)
			if err != nil {
				return
			} else {
				hasMoved = true
				break
			}
		}
		if !hasMoved {
			return
		}
	}

	return
}

func moveAwayOneNormalVolume(commandEnv *CommandEnv, volumeReplicas map[uint32][]*VolumeReplica, vol *master_pb.VolumeInformationMessage, thisNode *Node, otherNodes []*Node, applyChange bool) (hasMoved bool, err error) {
	sort.Slice(otherNodes, func(i, j int) bool {
		return otherNodes[i].localVolumeRatio() < otherNodes[j].localVolumeRatio()
	})

	for i := 0; i < len(otherNodes); i++ {
		emptyNode := otherNodes[i]
		hasMoved, err = maybeMoveOneVolume(commandEnv, volumeReplicas, thisNode, vol, emptyNode, applyChange)
		if err != nil {
			return
		}
		if hasMoved {
			break
		}
	}
	return
}

func nodesOtherThan(volumeServers []*Node, thisServer string) (thisNode *Node, otherNodes []*Node) {
	for _, node := range volumeServers {
		if node.info.Id == thisServer {
			thisNode = node
			continue
		}
		otherNodes = append(otherNodes, node)
	}
	return
}

func ecNodesOtherThan(volumeServers []*EcNode, thisServer string) (thisNode *EcNode, otherNodes []*EcNode) {
	for _, node := range volumeServers {
		if node.info.Id == thisServer {
			thisNode = node
			continue
		}
		otherNodes = append(otherNodes, node)
	}
	return
}
