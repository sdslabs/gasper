package python

import (
	"fmt"

	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
)

func startServer(runCommand string, env *types.ApplicationEnv) (string, types.ResponseError) {
	cmd := []string{
		"bash", "-c",
		fmt.Sprintf("source venv/bin/activate && %s", runCommand),
	}
	execID, err := docker.ExecDetachedProcess(env.Context, env.Client, env.ContainerID, cmd)
	if err != nil {
		return execID, types.NewResErr(500, "failed to start the server", err)
	}
	return execID, nil
}

func installRequirements(path string, env *types.ApplicationEnv) (string, types.ResponseError) {
	cmd := []string{
		"bash", "-c",
		fmt.Sprintf("source venv/bin/activate && pip install -r %s", path),
	}
	execID, err := docker.ExecDetachedProcess(env.Context, env.Client, env.ContainerID, cmd)
	if err != nil {
		return execID, types.NewResErr(500, "failed to install requirements", err)
	}
	return execID, nil
}

func createVenv(env *types.ApplicationEnv, version string) (string, types.ResponseError) {
	var path string
	switch version {
	case "python2":
		path = "/usr/bin/python2"
	case "python3":
		path = "/usr/bin/python3"
	default:
		return "", types.NewResErr(400, "accepted values of python_version are 'python2' or 'python3'", nil)
	}
	cmd := []string{"virtualenv", "-p", path, "venv"}
	execID, err := docker.ExecDetachedProcess(env.Context, env.Client, env.ContainerID, cmd)
	if err != nil {
		return execID, types.NewResErr(500, "failed to create python virtual environment", err)
	}
	return execID, nil
}
