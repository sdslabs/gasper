package docker

import (
	dockerTypes "github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// InspectContainerState returns the state of the container using the containerID
func InspectContainerState(containerID string) (*dockerTypes.ContainerState, error) {
	ctx := context.Background()
	containerStatus, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}
	return containerStatus.ContainerJSONBase.State, nil
}
