package api

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdslabs/SWS/lib/git"

	"github.com/docker/docker/client"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
	"golang.org/x/net/context"
)

// CreateBasicApplication spawns a new container with the application of a particular service
func CreateBasicApplication(name, location, url, httpPort, sshPort string, appConf *types.ApplicationConfig) (string, *types.ResponseError) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", types.NewResponseError(500, "", err)
	}
	var (
		storepath, _ = os.Getwd()
		confFileName = fmt.Sprintf("%s.sws.conf", name)
		workdir      = fmt.Sprintf("/SWS/%s", name)
		storedir     = filepath.Join(storepath, fmt.Sprintf("storage/%s", name))
	)

	// Step 1: clone the repo in the storage
	err = os.MkdirAll(storedir, 0755)
	// TODO: check if does exist -> name of the project changed and other issues
	// Check if dir available...?
	if err != nil {
		return "", types.NewResponseError(500, "failed to create storage directory [Mkdir]", err)
	}
	err = git.Clone(url, storedir)
	if err != nil {
		return "", types.NewResponseError(500, "failed to clone application [Clone]", err)
	}

	// Step 2: create the container
	containerID, err := docker.CreateContainer(ctx, cli, appConf.DockerImage, httpPort, sshPort, workdir, storedir, name)
	if err != nil {
		return "", types.NewResponseError(500, "failed to create new container [CreateContainer]", err)
	}
	// Step 3: write config to the container
	confFile := []byte(appConf.ConfFunction(name, location))
	archive, err := utils.NewTarArchiveFromContent(confFile, confFileName, 0644)
	if err != nil {
		return "", types.NewResponseError(500, "failed to write conf file [NewTarArchiveFromContent]", err)
	}
	err = docker.CopyToContainer(ctx, cli, containerID, "/etc/nginx/conf.d/", archive)
	if err != nil {
		return "", types.NewResponseError(500, "failed to write conf file [CopyToContainer]", err)
	}
	// Step 4: start the container
	err = docker.StartContainer(ctx, cli, containerID)
	if err != nil {
		return "", types.NewResponseError(500, "failed to start container [StartContainer]", err)
	}
	return containerID, nil
}
