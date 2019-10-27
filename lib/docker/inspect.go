package docker

import (
	"encoding/json"

	"github.com/sdslabs/gasper/types"
	"golang.org/x/net/context"
)

// InspectContainerState returns the state of the container using the containerID
func InspectContainerState(containerID string) (types.M, error) {
	ctx := context.Background()
	containerStatus, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}

	var containerStatusInterface types.M
	marshalledInterface, err := json.Marshal(containerStatus.ContainerJSONBase.State)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(marshalledInterface, &containerStatusInterface)

	return containerStatusInterface, nil
}
