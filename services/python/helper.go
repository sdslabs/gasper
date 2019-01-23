package python

import (
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
)

func startServer(index string, env *types.ApplicationEnv) (string, types.ResponseError) {
	cmd := []string{"python", index}
	execID, err := docker.ExecDetachedProcess(env.Context, env.Client, env.ContainerID, cmd)
	if err != nil {
		return execID, types.NewResErr(500, "failed to start the server", err)
	}
	return execID, nil
}

func installRequirements(path string, env *types.ApplicationEnv) (string, types.ResponseError) {
	cmd := []string{"pip", "install", "-r", path}
	execID, err := docker.ExecDetachedProcess(env.Context, env.Client, env.ContainerID, cmd)
	if err != nil {
		return execID, types.NewResErr(500, "failed to install requirements", err)
	}
	return execID, nil
}
