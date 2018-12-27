package php

import (
	"github.com/docker/docker/client"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
	"golang.org/x/net/context"
)

// installPackages installs dependancies for the specific microservice
func installPackages(path, containerID string) (string, *types.ResponseError) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", types.NewResponseError(500, "", err)
	}
	cmd := []string{"composer", "install", "-d", path}
	execID, err := docker.ExecDetachedProcess(ctx, cli, containerID, cmd)
	if err != nil {
		return "", types.NewResponseError(500, "Failed to perform composer install in the container", err)
	}
	return execID, nil
}
