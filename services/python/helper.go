package python

import (
	"fmt"
	"strings"

	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
)

func startServer(index string, args []string, env *types.ApplicationEnv) (string, types.ResponseError) {
	arguments := strings.Join(args, " ")
	serveCmd := fmt.Sprintf(`python %s %s &> /proc/1/fd/1`, index, arguments)
	cmd := []string{"bash", "-c", serveCmd}
	execID, err := docker.ExecDetachedProcess(env.Context, env.Client, env.ContainerID, cmd)
	if err != nil {
		return execID, types.NewResErr(500, "failed to start the server", err)
	}
	return execID, nil
}

func installRequirements(path string, env *types.ApplicationEnv) (string, types.ResponseError) {
	cmd := []string{"bash", "-c", fmt.Sprintf(`pip install -r %s &> /proc/1/fd/1`, path)}
	execID, err := docker.ExecDetachedProcess(env.Context, env.Client, env.ContainerID, cmd)
	if err != nil {
		return execID, types.NewResErr(500, "failed to install requirements", err)
	}
	return execID, nil
}
