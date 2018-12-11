package api

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
	"golang.org/x/net/context"
)

// CreateBasicApplication spawns a new container with the application of a particular service
func CreateBasicApplication(name, port80, port22 string, appConf *types.ApplicationConfig) *types.ResponseError {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return types.NewResponseError(500, "", err)
	}
	var (
		workDir      = fmt.Sprintf("/SWS/%s", name)
		confFileName = fmt.Sprintf("%s.sws.conf", name)
	)

	// Step 1: create the container
	containerID, err := docker.CreateContainer(ctx, cli, appConf.DockerImage, port80, port22, workDir, name)
	if err != nil {
		return types.NewResponseError(500, "failed to create new container [CreateContainer]", err)
	}
	// Step 2: write config to the container
	confFile := []byte(appConf.ConfFunction(name))
	archive, err := utils.NewTarArchiveFromContent(confFile, confFileName, 0644)
	if err != nil {
		return types.NewResponseError(500, "failed to write conf file [NewTarArchiveFromContent]", err)
	}
	err = docker.CopyToContainer(ctx, cli, containerID, "/etc/nginx/conf.d/", archive)
	if err != nil {
		return types.NewResponseError(500, "failed to write conf file [CopyToContainer]", err)
	}
	// Step 3: clone the repo and add it to the container
	// TODO: Step 3
	// Step 4: start the container
	err = docker.StartContainer(ctx, cli, containerID)
	if err != nil {
		return types.NewResponseError(500, "failed to start container [StartContainer]", err)
	}
	return nil
}
