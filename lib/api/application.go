package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/git"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
	gogit "gopkg.in/src-d/go-git.v4"
)

func cloneRepo(url, storedir, accessToken string, mutex map[string]chan types.ResponseError) {
	err := os.MkdirAll(storedir, 0755)
	if err != nil {
		mutex["clone"] <- types.NewResErr(500, "storage directory not created", err)
		return
	}

	if accessToken == "" {
		err = git.Clone(url, storedir)
	} else {
		err = git.CloneWithToken(url, storedir, accessToken)
	}

	if err != nil {
		if err == gogit.ErrRepositoryAlreadyExists {
			mutex["clone"] <- types.NewResErr(500, "repository already exists", err)
		} else {
			mutex["clone"] <- types.NewResErr(500, "repository not cloned", err)
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
	httpPort string,
	env map[string]interface{},
	appContext map[string]interface{},
	appConf *types.ApplicationConfig,
	mutex map[string]chan types.ResponseError) {

	var err error
	// create the container
	appEnv.ContainerID, err = docker.CreateContainer(appEnv.Context, appEnv.Client, appConf.DockerImage, httpPort, workdir, storedir, name, env)
	if err != nil {
		mutex["setup"] <- types.NewResErr(500, "container not created", err)
		return
	}

	// write config to the container
	confFile := []byte(appConf.ConfFunction(name, appContext))
	archive, err := utils.NewTarArchiveFromContent(confFile, confFileName, 0644)
	if err != nil {
		mutex["setup"] <- types.NewResErr(500, "container conf file not written", err)
		return
	}
	err = docker.CopyToContainer(appEnv.Context, appEnv.Client, appEnv.ContainerID, "/etc/nginx/conf.d/", archive)
	if err != nil {
		mutex["setup"] <- types.NewResErr(500, "container conf file not written", err)
		return
	}

	// start the container
	err = docker.StartContainer(appEnv.Context, appEnv.Client, appEnv.ContainerID)
	if err != nil {
		mutex["setup"] <- types.NewResErr(500, "container not started", err)
		return
	}
	mutex["setup"] <- nil
}

// CreateBasicApplication spawns a new container with the application of a particular service
func CreateBasicApplication(name, url, accessToken, httpPort string, env, appContext map[string]interface{}, appConf *types.ApplicationConfig) (*types.ApplicationEnv, []types.ResponseError) {
	appEnv, err := types.NewAppEnv()
	if err != nil {
		return nil, []types.ResponseError{types.NewResErr(500, "", err), nil}
	}

	var (
		storepath, _ = os.Getwd()
		confFileName = fmt.Sprintf("%s.sws.conf", name)
		workdir      = fmt.Sprintf("%s/%s", configs.SWSConfig["projectRoot"].(string), name)
		storedir     = filepath.Join(storepath, fmt.Sprintf("storage/%s", name))
	)

	var mutex = map[string]chan types.ResponseError{
		"setup": make(chan types.ResponseError),
		"clone": make(chan types.ResponseError),
	}

	// Step 1: clone the repo in the storage
	go cloneRepo(url, storedir, accessToken, mutex)

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
		env,
		appContext,
		appConf,
		mutex)

	setupErr := <-mutex["setup"]
	cloneErr := <-mutex["clone"]

	setupFlag := false
	cloneFlag := false

	if cloneErr != nil {
		if cloneErr.Message() != "repository already exists" {
			cloneFlag = true
		}
	}

	if setupErr != nil {
		if setupErr.Message() != "container not created" {
			setupFlag = true
		}
	}

	if setupFlag || cloneFlag {
		go utils.FullCleanup(name)
	}

	return appEnv, []types.ResponseError{setupErr, cloneErr}
}

// SetupApplication sets up a basic container for the application with all the prerequisites
func SetupApplication(appConf *types.ApplicationConfig, data map[string]interface{}) (*types.ApplicationEnv, types.ResponseError) {
	ports, err := utils.GetFreePorts(1)
	if err != nil {
		return nil, types.NewResErr(500, "free ports unavailable", err)
	}
	if len(ports) < 1 {
		return nil, types.NewResErr(500, "not enough free ports available", nil)
	}
	httpPort := ports[0]

	var env map[string]interface{}

	if data["env"] != nil {
		env = data["env"].(map[string]interface{})
	}

	accessToken := ""

	if data["git_access_token"] != nil {
		accessToken = data["git_access_token"].(string)
	}

	appEnv, errList := CreateBasicApplication(
		data["name"].(string),
		data["url"].(string),
		accessToken,
		strconv.Itoa(httpPort),
		env,
		data["context"].(map[string]interface{}),
		appConf)

	for _, e := range errList {
		if e != nil {
			return nil, e
		}
	}

	runCommands := false
	rcFile := data["context"].(map[string]interface{})["rcFile"]
	if rcFile != nil {
		runCommands = rcFile.(bool)
	}
	data["context"].(map[string]interface{})["rcFile"] = runCommands

	if runCommands {
		_, err = docker.ExecDetachedProcess(
			appEnv.Context,
			appEnv.Client,
			appEnv.ContainerID,
			[]string{"/bin/bash", configs.SWSConfig["rcFile"].(string)})
		if err != nil {
			// this error cannot be ignored; the chances of error here are very less
			// but if an error arises, this means there's some issue with "execing"
			// any process in the container => there's a problem with the container
			// hence we also run the cleanup here so that nothing else goes wrong
			go utils.FullCleanup(data["name"].(string))
			return nil, types.NewResErr(500, "cannot exec rc file", err)
		}
	}

	data["httpPort"] = httpPort
	data["containerID"] = appEnv.ContainerID
	data["hostIP"] = utils.HostIP

	return appEnv, nil
}
