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
	contStatus, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}

	var contStatusInterface map[string]interface{}
	marshalledInterface, err := json.Marshal(contStatus.ContainerJSONBase.State)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(marshalledInterface, &contStatusInterface)

	return contStatusInterface, nil
}
