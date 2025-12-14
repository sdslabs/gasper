package docker

import (
	"bytes"
	"errors"
	"io"
	"strconv"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ExecDetachedProcess executes a command in detached form, returns the id of the process
// Command of the exec format: mkdir folder => ["mkdir", "folder"]
func ExecDetachedProcess(containerID string, command []string) (string, error) {
	// TODO: check if container is up and running first
	ctx := context.Background()
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

// ExecProcess executes a command in a blocing manner and returns the id of the process
func ExecProcess(containerID string, command []string) (string, error) {
	ctx := context.Background()
	config := types.ExecConfig{
		Detach:       false,
		Tty:          true,
		Cmd:          command,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
	}
	execProcess, err := cli.ContainerExecCreate(ctx, containerID, config)
	if err != nil {
		return "", err
	}
	execID := execProcess.ID
	if execID == "" {
		return "", errors.New("empty exec ID")
	}
	err = cli.ContainerExecStart(ctx, execID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		return "", err
	}
	return execID, nil
}

// ExecProcessWthStream executes a command in a blocing manner and returns output stream as string
func ExecProcessWthStream(containerID string, command []string) (string, error) {
	ctx := context.Background()
	config := types.ExecConfig{
		Detach:       false,
		Tty:          true,
		Cmd:          command,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
	}	
	execProcess, err := cli.ContainerExecCreate(ctx, containerID, config)
	if err != nil {
		return "", err
	}
	execID := execProcess.ID
	if execID == "" {
		return "", errors.New("empty exec ID")
	}

	resp, err := cli.ContainerExecAttach(ctx, execID, types.ExecConfig{Detach: false, Tty: true})
	if err != nil {
		return "", err
	}
	defer resp.Close()

	var outputBuffer bytes.Buffer
	_, err = io.Copy(&outputBuffer, resp.Reader)
	if err != nil {
		return "", err
	}

	statusCode, err := cli.ContainerExecInspect(ctx, execID)
	if err != nil {
		return "", err
	}

	if statusCode.ExitCode != 0 {
		return "", errors.New("command execution failed with exit code " + strconv.Itoa(statusCode.ExitCode))
	}

	return outputBuffer.String(), nil
}	
	