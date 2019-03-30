package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// DeleteContainer deletes a docker container
func DeleteContainer(containerID string) error {
	ctx := context.Background()
	cli, _ := client.NewEnvClient()
	err := StopContainer(ctx, cli, containerID)
	
	if err != nil {
		return err
	}
	
	err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})

	if err != nil {
		return err
	}

	return nil
}
