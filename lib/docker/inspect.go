package docker

import (
	"encoding/json"

	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// InspectContainerState returns the state of the container using the containerID
func InspectContainerState(containerID string) (map[string]interface{}, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	containerStatus, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}

	var containerStatusInterface map[string]interface{}
	marshalledInterface, err := json.Marshal(containerStatus.ContainerJSONBase.State)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(marshalledInterface, &containerStatusInterface)

	return containerStatusInterface, nil
}
