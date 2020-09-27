package shell

import (
	"context"
	"flag"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"

	"github.com/chrislusf/seaweedfs/weed/operation"
	"github.com/chrislusf/seaweedfs/weed/pb/volume_server_pb"
	"github.com/chrislusf/seaweedfs/weed/storage/needle"
)

func init() {
	Commands = append(Commands, &commandVolumeTierUpload{})
}

type commandVolumeTierUpload struct {
}

func (c *commandVolumeTierUpload) Name() string {
	return "volume.tier.upload"
}

func (c *commandVolumeTierUpload) Help() string {
	return `upload the dat file of a volume to a remote tier

	volume.tier.upload [-collection=""] [-fullPercent=95] [-quietFor=1h]
	volume.tier.upload [-collection=""] -volumeId=<volume_id> -dest=<storage_backend> [-keepLocalDatFile]

	e.g.:
	volume.tier.upload -volumeId=7 -dest=s3
	volume.tier.upload -volumeId=7 -dest=s3.default

	The <storage_backend> is defined in master.toml.
	For example, "s3.default" in [storage.backend.s3.default]

	This command will move the dat file of a volume to a remote tier.

	SeaweedFS enables scalable and fast local access to lots of files, 
	and the cloud storage is slower by cost efficient. How to combine them together?

	Usually the data follows 80/20 rule: only 20% of data is frequently accessed.
	We can offload the old volumes to the cloud.

	With this, SeaweedFS can be both fast and scalable, and infinite storage space.
	Just add more local SeaweedFS volume servers to increase the throughput.

	The index file is still local, and the same O(1) disk read is applied to the remote file.

`
}

func (c *commandVolumeTierUpload) Do(args []string, commandEnv *CommandEnv, writer io.Writer) (err error) {

	if err = commandEnv.confirmIsLocked(); err != nil {
		return
	}

	tierCommand := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	volumeId := tierCommand.Int("volumeId", 0, "the volume id")
	collection := tierCommand.String("collection", "", "the collection name")
	fullPercentage := tierCommand.Float64("fullPercent", 95, "the volume reaches the percentage of max volume size")
	quietPeriod := tierCommand.Duration("quietFor", 24*time.Hour, "select volumes without no writes for this period")
	dest := tierCommand.String("dest", "", "the target tier name")
	keepLocalDatFile := tierCommand.Bool("keepLocalDatFile", false, "whether keep local dat file")
	if err = tierCommand.Parse(args); err != nil {
		return nil
	}

	vid := needle.VolumeId(*volumeId)

	// volumeId is provided
	if vid != 0 {
		return doVolumeTierUpload(commandEnv, writer, *collection, vid, *dest, *keepLocalDatFile)
	}

	// apply to all volumes in the collection
	// reusing collectVolumeIdsForEcEncode for now
	volumeIds, err := collectVolumeIdsForEcEncode(commandEnv, *collection, *fullPercentage, *quietPeriod)
	if err != nil {
		return err
	}
	fmt.Printf("tier upload volumes: %v\n", volumeIds)
	for _, vid := range volumeIds {
		if err = doVolumeTierUpload(commandEnv, writer, *collection, vid, *dest, *keepLocalDatFile); err != nil {
			return err
		}
	}

	return nil
}

func doVolumeTierUpload(commandEnv *CommandEnv, writer io.Writer, collection string, vid needle.VolumeId, dest string, keepLocalDatFile bool) (err error) {
	// find volume location
	locations, found := commandEnv.MasterClient.GetLocations(uint32(vid))
	if !found {
		return fmt.Errorf("volume %d not found", vid)
	}

	err = markVolumeReadonly(commandEnv.option.GrpcDialOption, needle.VolumeId(vid), locations)
	if err != nil {
		return fmt.Errorf("mark volume %d as readonly on %s: %v", vid, locations[0].Url, err)
	}

	// copy the .dat file to remote tier
	err = uploadDatToRemoteTier(commandEnv.option.GrpcDialOption, writer, needle.VolumeId(vid), collection, locations[0].Url, dest, keepLocalDatFile)
	if err != nil {
		return fmt.Errorf("copy dat file for volume %d on %s to %s: %v", vid, locations[0].Url, dest, err)
	}

	return nil
}

func uploadDatToRemoteTier(grpcDialOption grpc.DialOption, writer io.Writer, volumeId needle.VolumeId, collection string, sourceVolumeServer string, dest string, keepLocalDatFile bool) error {

	err := operation.WithVolumeServerClient(sourceVolumeServer, grpcDialOption, func(volumeServerClient volume_server_pb.VolumeServerClient) error {
		stream, copyErr := volumeServerClient.VolumeTierMoveDatToRemote(context.Background(), &volume_server_pb.VolumeTierMoveDatToRemoteRequest{
			VolumeId:               uint32(volumeId),
			Collection:             collection,
			DestinationBackendName: dest,
			KeepLocalDatFile:       keepLocalDatFile,
		})

		var lastProcessed int64
		for {
			resp, recvErr := stream.Recv()
			if recvErr != nil {
				if recvErr == io.EOF {
					break
				} else {
					return recvErr
				}
			}

			processingSpeed := float64(resp.Processed-lastProcessed) / 1024.0 / 1024.0

			fmt.Fprintf(writer, "copied %.2f%%, %d bytes, %.2fMB/s\n", resp.ProcessedPercentage, resp.Processed, processingSpeed)

			lastProcessed = resp.Processed
		}

		return copyErr
	})

	return err

}
