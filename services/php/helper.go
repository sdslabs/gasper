package php

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
	"golang.org/x/net/context"
)

func installPackages(path string) (string, *types.ResponseError) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", types.NewResponseError(500, "", err)
	}
	var composerInstallCmd string = "cd " + path + " && composer install"
	cmd := []string{"bash", "-c", composerInstallCmd}
	execId, err := docker.ExecDetachedProcess(ctx, cli, "c859d83d6ac1", cmd)
	fmt.Println(execId)
	if err != nil {
		return "", types.NewResponseError(500, "Failed to perform composer install in the container", err)
	}
	return execId, nil
}
