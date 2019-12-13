package api

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/commons"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/git"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
	gogit "gopkg.in/src-d/go-git.v4"
)

func cloneRepo(app types.Application, storedir string, mutex map[string]chan types.ResponseError) {
	err := os.MkdirAll(storedir, 0755)
	if err != nil {
		mutex["clone"] <- types.NewResErr(500, "storage directory not created", err)
		return
	}

	if app.HasGitAccessToken() {
		err = git.CloneWithToken(
			app.GetGitRepositoryURL(),
			app.GetGitRepositoryBranch(),
			storedir,
			app.GetGitAccessToken(),
		)
	} else {
		err = git.Clone(
			app.GetGitRepositoryURL(),
			app.GetGitRepositoryBranch(),
			storedir,
		)
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

func setupContainer(app types.Application, storedir string, mutex map[string]chan types.ResponseError) {
	var (
		confFileName = fmt.Sprintf("%s.gasper.conf", app.GetName())
		workdir      = fmt.Sprintf("%s/%s", configs.GasperConfig.ProjectRoot, app.GetName())
	)

	// create the container
	containerID, err := docker.CreateContainer(&types.ApplicationContainer{
		Name:            app.GetName(),
		Image:           app.GetDockerImage(),
		ApplicationPort: app.GetApplicationPort(),
		ContainerPort:   app.GetContainerPort(),
		WorkDir:         workdir,
		StoreDir:        storedir,
		Env:             app.GetEnvVars(),
		Memory:          app.GetMemoryLimit(),
		CPU:             app.GetCPULimit(),
		NameServers:     app.GetNameServers(),
	})

	if err != nil {
		mutex["setup"] <- types.NewResErr(500, "container not created", err)
		return
	}

	app.SetContainerID(containerID)

	// For PHP and Static applications, a nginx configuration is necessary
	if app.HasConfGenerator() {
		// write config to the container
		confFile := []byte(app.InvokeConfGenerator(app.GetName(), app.GetIndex()))
		archive, err := utils.NewTarArchiveFromContent(confFile, confFileName, 0644)
		if err != nil {
			mutex["setup"] <- types.NewResErr(500, "container conf file not written", err)
			return
		}
		err = docker.CopyToContainer(app.GetContainerID(), "/etc/nginx/conf.d/", archive)
		if err != nil {
			mutex["setup"] <- types.NewResErr(500, "container conf file not written", err)
			return
		}
	}

	// start the container
	err = docker.StartContainer(app.GetContainerID())
	if err != nil {
		mutex["setup"] <- types.NewResErr(500, "container not started", err)
		return
	}
	mutex["setup"] <- nil
}

// CreateBasicApplication spawns a new container with the application of a particular service
func CreateBasicApplication(app types.Application) []types.ResponseError {
	var (
		storepath, _ = os.Getwd()
		storedir     = filepath.Join(storepath, fmt.Sprintf("storage/%s", app.GetName()))
	)

	var mutex = map[string]chan types.ResponseError{
		"setup": make(chan types.ResponseError),
		"clone": make(chan types.ResponseError),
	}

	// Step 1: clone the repo in the storage
	go cloneRepo(app, storedir, mutex)

	// Step 2: setup the container
	go setupContainer(app, storedir, mutex)

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
		go commons.AppFullCleanup(app.GetName())
	}

	return []types.ResponseError{setupErr, cloneErr}
}

// SetupApplication sets up a basic container for the application with all the prerequisites
func SetupApplication(app types.Application) types.ResponseError {
	containerPort, err := utils.GetFreePort()
	if err != nil {
		return types.NewResErr(500, "No free port available", err)
	}

	app.SetContainerPort(containerPort)

	errList := CreateBasicApplication(app)

	for _, e := range errList {
		if e != nil {
			return e
		}
	}

	if app.HasRcFile() {
		cmd := []string{"sh", "-c",
			fmt.Sprintf(`chmod 755 ./%s &> /proc/1/fd/1 && ./%s &> /proc/1/fd/1`,
				configs.GasperConfig.RcFile, configs.GasperConfig.RcFile)}

		_, err = docker.ExecDetachedProcess(app.GetContainerID(), cmd)
		if err != nil {
			// this error cannot be ignored; the chances of error here are very less
			// but if an error arises, this means there's some issue with "execing"
			// any process in the container => there's a problem with the container
			// hence we also run the cleanup here so that nothing else goes wrong
			go commons.AppFullCleanup(app.GetName())
			return types.NewResErr(500, "cannot exec rc file", err)
		}
	} else {
		go buildAndRun(app)
	}

	return nil
}
