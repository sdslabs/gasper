package docker

import (
	dockerTypes "github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

const (
	// Strings for ContainterHealth
	Container_Healthy = "healthy"
	Container_Unhealthy = "unhealthy"
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

// ContainerHealth returns the health status of the container
func InspectContainerHealth(containerID string) (string, error) {
	ctx := context.Background()
	health, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", err
	}
	return health.State.Health.Status, nil
}
