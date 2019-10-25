package docker

import (
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// DeleteContainer deletes a docker container
func DeleteContainer(containerID string) error {
	ctx := context.Background()
	err := StopContainer(containerID)

	if err != nil {
		return err
	}

	err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})

	if err != nil {
		return err
	}

	return nil
}
