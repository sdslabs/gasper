package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

func setupContainer(app types.Application, storedir string, setup chan types.ResponseError) {
	confFileName := fmt.Sprintf("%s.gasper.conf", app.GetName())
	workdir := fmt.Sprintf("%s/%s", configs.GasperConfig.ProjectRoot, app.GetName())

	// create the container
	containerID, err := docker.CreateApplicationContainer(types.ApplicationContainer{
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
		if strings.HasPrefix(err.Error(), "Error response from daemon: Conflict.") {
			setup <- types.NewResErr(500, "container already exists", err)
		} else {
			setup <- types.NewResErr(500, "container not created", err)
		}
		return
	}

	app.SetContainerID(containerID)

	// For PHP and Static applications, a nginx configuration is necessary
	if app.HasConfGenerator() {
		// write config to the container
		confFile := []byte(app.InvokeConfGenerator(app.GetName(), app.GetIndex()))
		archive, err := utils.NewTarArchiveFromContent(confFile, confFileName, 0644)
		if err != nil {
			setup <- types.NewResErr(500, "container conf file not written", err)
			return
		}
		err = docker.CopyToContainer(app.GetContainerID(), "/etc/nginx/conf.d/", archive)
		if err != nil {
			setup <- types.NewResErr(500, "container conf file not written", err)
			return
		}
	}

	// start the container
	err = docker.StartContainer(app.GetContainerID())
	if err != nil {
		setup <- types.NewResErr(500, "container not started", err)
		return
	}
	setup <- nil
}

// createBasicApplication spawns a new container with the application of a particular service
func CreateBasicApplication(app types.Application) []types.ResponseError {
	storepath, _ := os.Getwd()
	storedir := filepath.Join(storepath, fmt.Sprintf("storage/%s", app.GetName()))
	setup := make(chan types.ResponseError)

	// Step 1: setup the container
	go setupContainer(app, storedir, setup)

	return []types.ResponseError{<-setup}
}

// SetupApplication sets up a basic container for the application with all the prerequisites
func SetupApplication(app types.Application) types.ResponseError {
	containerPort, err := utils.GetFreePort()
	if err != nil {
		return types.NewResErr(500, "No free port available", err)
	}

	app.SetContainerPort(containerPort)

	errList := CreateBasicApplication(app)

	for _, err := range errList {
		if err != nil {
			return err
		}
	}

	_, err = docker.ExecProcess(app.GetContainerID(), []string{"git", "init"})
	if err != nil {
		return types.NewResErr(500, "git init unsuccessful", err)
	}

	var cloneURL string
	if len(app.GetGitAccessToken()) > 0 {
		split := strings.Split(app.GetGitRepositoryURL(), "//")
		cloneURL = fmt.Sprintf("https://oauth2:%s@%s", app.GetGitAccessToken() , split[1])
	} else {
		cloneURL = app.GetGitRepositoryURL()
	}

	_, err = docker.ExecProcess(app.GetContainerID(), []string{"git", "remote", "add", "origin", cloneURL})
	if err != nil {
		return types.NewResErr(500, "setting remote unsuccessful", err)
	}

	_, err = docker.ExecProcess(app.GetContainerID(), []string{"git", "pull", "origin", app.GetGitRepositoryBranch()})
	if err != nil {
		return types.NewResErr(500, "pulling contents unsuccessful", err)
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
			return types.NewResErr(500, "cannot exec rc file", err)
		}
	} else {
		go buildAndRun(app)
	}

	return nil
}
