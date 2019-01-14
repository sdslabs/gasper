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
func CreateBasicApplication(name, contextParam, url, httpPort, sshPort string, appConf *types.ApplicationConfig) (*types.ApplicationEnv, types.ResponseError) {
	appEnv, err := types.NewAppEnv()
	if err != nil {
		return nil, types.NewResErr(500, "", err)
	}

	var (
		storepath, _ = os.Getwd()
		confFileName = fmt.Sprintf("%s.sws.conf", name)
		workdir      = fmt.Sprintf("/SWS/%s", name)
		storedir     = filepath.Join(storepath, fmt.Sprintf("storage/%s", name))
	)

	// Step 1: clone the repo in the storage
	err = os.MkdirAll(storedir, 0755)
	if err != nil {
		return nil, types.NewResErr(500, "storage directory not created", err)
	}
	err = git.Clone(url, storedir)
	if err != nil {
		return nil, types.NewResErr(500, "repo not cloned", err)
	}

	// Step 2: create the container
	appEnv.ContainerID, err = docker.CreateContainer(appEnv.Context, appEnv.Client, appConf.DockerImage, httpPort, sshPort, workdir, storedir, name)
	if err != nil {
		return nil, types.NewResErr(500, "container not created", err)
	}

	// Step 3: write config to the container
	confFile := []byte(appConf.ConfFunction(name, contextParam))
	archive, err := utils.NewTarArchiveFromContent(confFile, confFileName, 0644)
	if err != nil {
		return appEnv, types.NewResErr(500, "container conf file not written", err)
	}
	err = docker.CopyToContainer(appEnv.Context, appEnv.Client, appEnv.ContainerID, "/etc/nginx/conf.d/", archive)
	if err != nil {
		return appEnv, types.NewResErr(500, "container conf file not written", err)
	}

	// Step 4: start the container
	err = docker.StartContainer(appEnv.Context, appEnv.Client, appEnv.ContainerID)
	if err != nil {
		return appEnv, types.NewResErr(500, "container not started", err)
	}
	return appEnv, nil
}
