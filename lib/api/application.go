package api

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdslabs/SWS/lib/git"

	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
)

// CreateBasicApplication spawns a new container with the application of a particular service
func CreateBasicApplication(name, location, url, httpPort, sshPort string, appConf *types.ApplicationConfig) (*types.ApplicationEnv, *types.ResponseError) {
	appEnv, err := types.NewAppEnv()
	if err != nil {
		return nil, types.NewResponseError(500, "", err)
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
		return nil, types.NewResponseError(500, "failed to create storage directory [Mkdir]", err)
	}
	err = git.Clone(url, storedir)
	if err != nil {
		return nil, types.NewResponseError(500, "failed to clone application [Clone]", err)
	}

	// Step 2: create the container
	appEnv.ContainerID, err = docker.CreateContainer(appEnv.Context, appEnv.Client, appConf.DockerImage, httpPort, sshPort, workdir, storedir, name)
	if err != nil {
		return nil, types.NewResponseError(500, "failed to create new container [CreateContainer]", err)
	}
	// Step 3: write config to the container
	confFile := []byte(appConf.ConfFunction(name, location))
	archive, err := utils.NewTarArchiveFromContent(confFile, confFileName, 0644)
	if err != nil {
		return appEnv, types.NewResponseError(500, "failed to write conf file [NewTarArchiveFromContent]", err)
	}
	err = docker.CopyToContainer(appEnv.Context, appEnv.Client, appEnv.ContainerID, "/etc/nginx/conf.d/", archive)
	if err != nil {
		return appEnv, types.NewResponseError(500, "failed to write conf file [CopyToContainer]", err)
	}
	// Step 4: start the container
	err = docker.StartContainer(appEnv.Context, appEnv.Client, appEnv.ContainerID)
	if err != nil {
		return appEnv, types.NewResponseError(500, "failed to start container [StartContainer]", err)
	}
	return appEnv, nil
}
