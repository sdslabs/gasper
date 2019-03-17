package docker

import (
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// ExecDetachedProcess executes a command in detached form, returns the id of the process
// Command of the exec format: mkdir folder => ["mkdir", "folder"]
func ExecDetachedProcess(ctx context.Context, cli *client.Client, containerID string, command []string) (string, error) {
	// TODO: check if container is up and running first
	config := types.ExecConfig{
		Detach: true,
		Cmd:    command,
	}
	execProcess, err := cli.ContainerExecCreate(ctx, containerID, config)
	if err != nil {
		return "", err
	}
	execID := execProcess.ID
	if execID == "" {
		return "", errors.New("empty exec ID")
	}
	err = cli.ContainerExecStart(ctx, execID, types.ExecStartCheck{Detach: true})
	if err != nil {
		return "", err
	}
	return execID, nil
}
