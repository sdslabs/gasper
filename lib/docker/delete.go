package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	Constants "github.com/sdslabs/SWS/lib/utils/constants"
)

// DeleteContainer deletes a docker container
func DeleteContainer(containerID string) error {
	ctx := context.Background()
	cli, _ := client.NewEnvClient()
	err := StopContainer(ctx, cli, containerID)
	str := "Error response from daemon: No such container: " + containerID
	if err != nil && err.Error() != str {
		return err
	}
	err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})

	// Error := errors.New("Error response from daemon: No such container: ")
	Error := Constants.ConstErr

	if err != nil && err != Error {
		return err
	}
	return nil
}
