package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/git"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
	gogit "gopkg.in/src-d/go-git.v4"
)

func cloneRepo(url, storedir string, mutex map[string]chan types.ResponseError) {
	err := os.MkdirAll(storedir, 0755)
	if err != nil {
		mutex["clone"] <- types.NewResErr(500, "storage directory not created", err)
		return
	}
	err = git.Clone(url, storedir)
	if err != nil {
		mutex["clone"] <- types.NewResErr(500, "repository not cloned", err)
		if err != gogit.ErrRepositoryAlreadyExists {
			slice := strings.Split(storedir, "/")
			appName := slice[len(slice)-1]
			utils.FullCleanup(appName)
		}
		return
	}
	mutex["clone"] <- nil
}

func setupContainer(
	appEnv *types.ApplicationEnv,
	storePath,
	confFileName,
	workdir,
	storedir,
	name,
	url,
	httpPort,
	sshPort string,
	env map[string]interface{},
	appContext map[string]interface{},
	appConf *types.ApplicationConfig,
	mutex map[string]chan types.ResponseError) {

	var err error
	// create the container
	appEnv.ContainerID, err = docker.CreateContainer(appEnv.Context, appEnv.Client, appConf.DockerImage, httpPort, sshPort, workdir, storedir, name, env)
	if err != nil {
		// return nil, types.NewResErr(500, "container not created", err)
		mutex["setup"] <- types.NewResErr(500, "container not created", err)
		return
	}

	// write config to the container
	confFile := []byte(appConf.ConfFunction(name, appContext))
	archive, err := utils.NewTarArchiveFromContent(confFile, confFileName, 0644)
	if err != nil {
		mutex["setup"] <- types.NewResErr(500, "container conf file not written", err)
		utils.FullCleanup(name)
		return
	}
	err = docker.CopyToContainer(appEnv.Context, appEnv.Client, appEnv.ContainerID, "/etc/nginx/conf.d/", archive)
	if err != nil {
		mutex["setup"] <- types.NewResErr(500, "container conf file not written", err)
		utils.FullCleanup(name)
		return
	}

	// start the container
	err = docker.StartContainer(appEnv.Context, appEnv.Client, appEnv.ContainerID)
	if err != nil {
		mutex["setup"] <- types.NewResErr(500, "container not started", err)
		utils.FullCleanup(name)
		return
	}
	mutex["setup"] <- nil
}

// CreateBasicApplication spawns a new container with the application of a particular service
func CreateBasicApplication(name, url, httpPort, sshPort string, env map[string]interface{}, appContext map[string]interface{}, appConf *types.ApplicationConfig) (*types.ApplicationEnv, []types.ResponseError) {
	appEnv, err := types.NewAppEnv()
	if err != nil {
		return nil, []types.ResponseError{types.NewResErr(500, "", err), nil}
	}

	var (
		storepath, _ = os.Getwd()
		confFileName = fmt.Sprintf("%s.sws.conf", name)
		workdir      = fmt.Sprintf("/SWS/%s", name)
		storedir     = filepath.Join(storepath, fmt.Sprintf("storage/%s", name))
	)

	var mutex = map[string]chan types.ResponseError{
		"setup": make(chan types.ResponseError),
		"clone": make(chan types.ResponseError),
	}

	// Step 1: clone the repo in the storage
	go cloneRepo(url, storedir, mutex)

	// Step 2: setup the container
	go setupContainer(
		appEnv,
		storepath,
		confFileName,
		workdir,
		storedir,
		name,
		url,
		httpPort,
		sshPort,
		env,
		appContext,
		appConf,
		mutex)

	return appEnv, []types.ResponseError{<-mutex["setup"], <-mutex["clone"]}
}
