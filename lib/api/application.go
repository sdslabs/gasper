package api

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/commons"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/git"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
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
	storePath, confFileName, workdir, storedir, name, url, httpPort, containerPort string,
	env, appContext types.M,
	appConf *types.ApplicationConfig,
	resources container.Resources,
	mutex map[string]chan types.ResponseError) {

	var err error
	// create the container
	appEnv.ContainerID, err = docker.CreateContainer(appConf.DockerImage, httpPort, containerPort, workdir, storedir, name, resources, env)
	if err != nil {
		mutex["setup"] <- types.NewResErr(500, "container not created", err)
		return
	}

	if appContext["port"] == nil {
		// write config to the container
		confFile := []byte(appConf.ConfFunction(name, appContext))
		archive, err := utils.NewTarArchiveFromContent(confFile, confFileName, 0644)
		if err != nil {
			mutex["setup"] <- types.NewResErr(500, "container conf file not written", err)
			return
		}
		err = docker.CopyToContainer(appEnv.ContainerID, "/etc/nginx/conf.d/", archive)
		if err != nil {
			mutex["setup"] <- types.NewResErr(500, "container conf file not written", err)
			return
		}
	}

	// start the container
	err = docker.StartContainer(appEnv.ContainerID)
	if err != nil {
		mutex["setup"] <- types.NewResErr(500, "container not started", err)
		return
	}
	mutex["setup"] <- nil
}

// CreateBasicApplication spawns a new container with the application of a particular service
func CreateBasicApplication(
	name, url, accessToken, httpPort, containerPort string,
	env, appContext types.M,
	appConf *types.ApplicationConfig,
	resources container.Resources) (*types.ApplicationEnv, []types.ResponseError) {

	appEnv, err := types.NewAppEnv()
	if err != nil {
		return nil, []types.ResponseError{types.NewResErr(500, "", err), nil}
	}

	var (
		storepath, _ = os.Getwd()
		confFileName = fmt.Sprintf("%s.gasper.conf", name)
		workdir      = fmt.Sprintf("%s/%s", configs.GasperConfig.ProjectRoot, name)
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
		containerPort,
		env,
		appContext,
		appConf,
		resources,
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
		go commons.AppFullCleanup(name)
	}

	return appEnv, []types.ResponseError{setupErr, cloneErr}
}

// SetupApplication sets up a basic container for the application with all the prerequisites
func SetupApplication(appConf *types.ApplicationConfig, data types.M) (*types.ApplicationEnv, types.ResponseError) {
	ports, err := utils.GetFreePorts(1)
	if err != nil {
		return nil, types.NewResErr(500, "free ports unavailable", err)
	}
	if len(ports) < 1 {
		return nil, types.NewResErr(500, "not enough free ports available", nil)
	}
	httpPort := ports[0]

	var env types.M

	if data["env"] != nil {
		env = data["env"].(types.M)
	}

	var resources container.Resources

	if data["resources"] == nil {
		data["resources"] = map[string]interface{}{
			"memory": docker.DefaultMemory,
			"cpu":    docker.DefaultCPUs,
		}
	}

	if data["resources"].(map[string]interface{})["memory"] != nil {
		resources.Memory = int64(data["resources"].(map[string]interface{})["memory"].(float64) * math.Pow(1024, 3))
	}
	if data["resources"].(map[string]interface{})["cpu"] != nil {
		resources.NanoCPUs = int64(data["resources"].(map[string]interface{})["cpu"].(float64) * math.Pow(10, 9))
	}

	accessToken := ""

	if data["git_access_token"] != nil {
		accessToken = data["git_access_token"].(string)
	}

	context := data["context"].(map[string]interface{})

	var containerPort string

	if context["port"] != nil {
		containerPort = context["port"].(string)
	} else {
		containerPort = "80"
	}

	appEnv, errList := CreateBasicApplication(
		data["name"].(string),
		data["url"].(string),
		accessToken,
		strconv.Itoa(httpPort),
		containerPort,
		env,
		data["context"].(map[string]interface{}),
		appConf,
		resources)

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
		cmd := []string{"sh", "-c", fmt.Sprintf(`chmod 755 ./%s &> /proc/1/fd/1 && ./%s &> /proc/1/fd/1`, configs.GasperConfig.RcFile, configs.GasperConfig.RcFile)}
		_, err = docker.ExecDetachedProcess(
			appEnv.ContainerID,
			cmd)
		if err != nil {
			// this error cannot be ignored; the chances of error here are very less
			// but if an error arises, this means there's some issue with "execing"
			// any process in the container => there's a problem with the container
			// hence we also run the cleanup here so that nothing else goes wrong
			go commons.AppFullCleanup(data["name"].(string))
			return nil, types.NewResErr(500, "cannot exec rc file", err)
		}
	}

	data["httpPort"] = httpPort
	data["containerID"] = appEnv.ContainerID
	data["hostIP"] = utils.HostIP

	return appEnv, nil
}
