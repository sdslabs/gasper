package node

import (
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
)

// installPackages function installs the dependancies for the app
func installPackages(appEnv *types.ApplicationEnv) (string, types.ResponseError) {
	cmd := []string{"bash", "-c", `export NVM_DIR="$HOME/.nvm" && [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"; npm install &> /proc/1/fd/1`}
	execID, err := docker.ExecDetachedProcess(appEnv.Context, appEnv.Client, appEnv.ContainerID, cmd)
	if err != nil {
		return "", types.NewResErr(500, "Failed to perform npm install in the container", err)
	}
	return execID, nil
}

// startApp function starts the app using pm2
func startApp(index string, appEnv *types.ApplicationEnv) (string, types.ResponseError) {
	cmd := []string{"bash", "-c", `export NVM_DIR="$HOME/.nvm" && [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"; pm2 start ` + index + ` &> /proc/1/fd/1`}
	execID, err := docker.ExecDetachedProcess(appEnv.Context, appEnv.Client, appEnv.ContainerID, cmd)
	if err != nil {
		return "", types.NewResErr(500, "Failed to perform start app in the container", err)
	}
	return execID, nil
}
