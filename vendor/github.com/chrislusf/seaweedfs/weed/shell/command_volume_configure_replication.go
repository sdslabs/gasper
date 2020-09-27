package shell

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/chrislusf/seaweedfs/weed/operation"
	"github.com/chrislusf/seaweedfs/weed/pb/master_pb"
	"github.com/chrislusf/seaweedfs/weed/pb/volume_server_pb"
	"github.com/chrislusf/seaweedfs/weed/storage/needle"
	"github.com/chrislusf/seaweedfs/weed/storage/super_block"
)

func init() {
	Commands = append(Commands, &commandVolumeConfigureReplication{})
}

type commandVolumeConfigureReplication struct {
}

func (c *commandVolumeConfigureReplication) Name() string {
	return "volume.configure.replication"
}

func (c *commandVolumeConfigureReplication) Help() string {
	return `change volume replication value

	This command changes a volume replication value. It should be followed by "volume.fix.replication".

`
}

func (c *commandVolumeConfigureReplication) Do(args []string, commandEnv *CommandEnv, writer io.Writer) (err error) {

	if err = commandEnv.confirmIsLocked(); err != nil {
		return
	}

	configureReplicationCommand := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	volumeIdInt := configureReplicationCommand.Int("volumeId", 0, "the volume id")
	replicationString := configureReplicationCommand.String("replication", "", "the intended replication value")
	if err = configureReplicationCommand.Parse(args); err != nil {
		return nil
	}

	if *replicationString == "" {
		return fmt.Errorf("empty replication value")
	}

	replicaPlacement, err := super_block.NewReplicaPlacementFromString(*replicationString)
	if err != nil {
		return fmt.Errorf("replication format: %v", err)
	}
	replicaPlacementInt32 := uint32(replicaPlacement.Byte())

	var resp *master_pb.VolumeListResponse
	err = commandEnv.MasterClient.WithClient(func(client master_pb.SeaweedClient) error {
		resp, err = client.VolumeList(context.Background(), &master_pb.VolumeListRequest{})
		return err
	})
	if err != nil {
		return err
	}

	vid := needle.VolumeId(*volumeIdInt)

	// find all data nodes with volumes that needs replication change
	var allLocations []location
	eachDataNode(resp.TopologyInfo, func(dc string, rack RackId, dn *master_pb.DataNodeInfo) {
		loc := newLocation(dc, string(rack), dn)
		for _, v := range dn.VolumeInfos {
			if v.Id == uint32(vid) && v.ReplicaPlacement != replicaPlacementInt32 {
				allLocations = append(allLocations, loc)
				continue
			}
		}
	})

	if len(allLocations) == 0 {
		return fmt.Errorf("no volume needs change")
	}

	for _, dst := range allLocations {
		err := operation.WithVolumeServerClient(dst.dataNode.Id, commandEnv.option.GrpcDialOption, func(volumeServerClient volume_server_pb.VolumeServerClient) error {
			resp, configureErr := volumeServerClient.VolumeConfigure(context.Background(), &volume_server_pb.VolumeConfigureRequest{
				VolumeId:    uint32(vid),
				Replication: replicaPlacement.String(),
			})
			if configureErr != nil {
				return configureErr
			}
			if resp.Error != "" {
				return errors.New(resp.Error)
			}
			return nil
		})

		if err != nil {
			return err
		}

	}

	return nil
}
