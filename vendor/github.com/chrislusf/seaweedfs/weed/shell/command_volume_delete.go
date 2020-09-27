package shell

import (
	"flag"
	"io"

	"github.com/chrislusf/seaweedfs/weed/storage/needle"
)

func init() {
	Commands = append(Commands, &commandVolumeDelete{})
}

type commandVolumeDelete struct {
}

func (c *commandVolumeDelete) Name() string {
	return "volume.delete"
}

func (c *commandVolumeDelete) Help() string {
	return `delete a live volume from one volume server

	volume.delete -node <volume server host:port> -volumeId <volume id>

	This command deletes a volume from one volume server.

`
}

func (c *commandVolumeDelete) Do(args []string, commandEnv *CommandEnv, writer io.Writer) (err error) {

	if err = commandEnv.confirmIsLocked(); err != nil {
		return
	}

	volDeleteCommand := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	volumeIdInt := volDeleteCommand.Int("volumeId", 0, "the volume id")
	nodeStr := volDeleteCommand.String("node", "", "the volume server <host>:<port>")
	if err = volDeleteCommand.Parse(args); err != nil {
		return nil
	}

	sourceVolumeServer := *nodeStr

	volumeId := needle.VolumeId(*volumeIdInt)

	return deleteVolume(commandEnv.option.GrpcDialOption, volumeId, sourceVolumeServer)

}
