package docker

import (
	// "encoding/json"

	"encoding/json"
	"io/ioutil"

	"github.com/sdslabs/gasper/types"
	"golang.org/x/net/context"
)

// ContainerStats returns container stats using the containerID
func ContainerStats(containerID string) (types.M, error) {
	ctx := context.Background()
	containerStats, err := cli.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(containerStats.Body)
	if err != nil {
		return nil, err
	}
	var containerStatsInterface types.M
	json.Unmarshal(body, &containerStatsInterface)
	return containerStatsInterface, nil
}
