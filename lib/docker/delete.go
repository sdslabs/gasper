package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// DeleteContainer deletes a docker container
func DeleteContainer(containerID string) {
	ctx := context.Background()
	cli, _ := client.NewEnvClient()
	StopContainer(ctx, cli, containerID)
	cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
}
