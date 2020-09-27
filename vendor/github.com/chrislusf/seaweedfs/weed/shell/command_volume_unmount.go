package shell

import (
	"context"
	"flag"
	"io"

	"github.com/chrislusf/seaweedfs/weed/operation"
	"github.com/chrislusf/seaweedfs/weed/pb/volume_server_pb"
	"github.com/chrislusf/seaweedfs/weed/storage/needle"
	"google.golang.org/grpc"
)

func init() {
	Commands = append(Commands, &commandVolumeUnmount{})
}

type commandVolumeUnmount struct {
}

func (c *commandVolumeUnmount) Name() string {
	return "volume.unmount"
}

func (c *commandVolumeUnmount) Help() string {
	return `unmount a volume from one volume server

	volume.unmount -node <volume server host:port> -volumeId <volume id>

	This command unmounts a volume from one volume server.

`
}

func (c *commandVolumeUnmount) Do(args []string, commandEnv *CommandEnv, writer io.Writer) (err error) {

	if err = commandEnv.confirmIsLocked(); err != nil {
		return
	}

	volUnmountCommand := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	volumeIdInt := volUnmountCommand.Int("volumeId", 0, "the volume id")
	nodeStr := volUnmountCommand.String("node", "", "the volume server <host>:<port>")
	if err = volUnmountCommand.Parse(args); err != nil {
		return nil
	}

	sourceVolumeServer := *nodeStr

	volumeId := needle.VolumeId(*volumeIdInt)

	return unmountVolume(commandEnv.option.GrpcDialOption, volumeId, sourceVolumeServer)

}

func unmountVolume(grpcDialOption grpc.DialOption, volumeId needle.VolumeId, sourceVolumeServer string) (err error) {
	return operation.WithVolumeServerClient(sourceVolumeServer, grpcDialOption, func(volumeServerClient volume_server_pb.VolumeServerClient) error {
		_, unmountErr := volumeServerClient.VolumeUnmount(context.Background(), &volume_server_pb.VolumeUnmountRequest{
			VolumeId: uint32(volumeId),
		})
		return unmountErr
	})
}
