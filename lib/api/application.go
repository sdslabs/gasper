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
// It is assumed that all the parameters have been checked and verified
// Returns containerID and ResponseError if any
func CreateBasicApplication(name, url, httpPort, sshPort string, appConf *types.ApplicationConfig) (string, *types.ResponseError) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", types.NewResponseError(500, "", err)
	}
	var (
		workDir      = "/SWS"
		confFileName = fmt.Sprintf("%s.sws.conf", name)
	)

	// Step 1: create the container
	containerID, err := docker.CreateContainer(ctx, cli, appConf.DockerImage, httpPort, sshPort, workDir, name)
	if err != nil {
		return "", types.NewResponseError(500, "failed to create new container [CreateContainer]", err)
	}

	// TODO: error from here should can create a garbage of containers so containers should be deleted
	// if there exists an error from this step

	// Step 2: write config to the container
	confFile := []byte(appConf.ConfFunction(name))
	archive, err := utils.NewTarArchiveFromContent(confFile, confFileName, 0644)
	if err != nil {
		return containerID, types.NewResponseError(500, "failed to write conf file [NewTarArchiveFromContent]", err)
	}
	err = docker.CopyToContainer(ctx, cli, containerID, "/etc/nginx/conf.d/", archive)
	if err != nil {
		return containerID, types.NewResponseError(500, "failed to write conf file [CopyToContainer]", err)
	}
	// Step 3: start the container
	err = docker.StartContainer(ctx, cli, containerID)
	if err != nil {
		return containerID, types.NewResponseError(500, "failed to start container [StartContainer]", err)
	}
	// Step 4: git clone application
	// App will be cloned in '/SWS/app' directory
	cloneCmd := utils.NewExecCmd([]string{fmt.Sprintf("git clone %s app", url)})
	_, err = docker.ExecDetachedProcess(ctx, cli, containerID, cloneCmd)
	if err != nil {
		return containerID, types.NewResponseError(500, "failed to clone the repository [ExecDetachedProcess]", err)
	}
	// ExecDetachedProcess also returns an execID which can be used to inspect the success of the clone process
	return containerID, nil
}
