package docker

import (
	"fmt"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// DeleteContainer deletes a docker container
func DeleteContainer(containerID string) error {
	ctx := context.Background()

	// Inspect the container to get its working directory
	containerJSON, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		utils.LogError("Docker-DeleteContainer-1", err)
		return err
	}
	workingDir := containerJSON.Config.WorkingDir
	if workingDir != "" {
		// Clear the working directory inside the container, including hidden files and directories
		cmd := []string{"sh", "-c", fmt.Sprintf("rm -rf %s/* %s/.*", workingDir, workingDir)}
		execConfig := types.ExecConfig{
			Cmd:          cmd,
			AttachStdout: true,
			AttachStderr: true,
			Privileged:   true,
		}

		execIDResp, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
		if err != nil {
			utils.LogError("Docker-DeleteContainer-2", err)
			return err
		}

		err = cli.ContainerExecStart(ctx, execIDResp.ID, types.ExecStartCheck{})
		if err != nil {
			utils.LogError("Docker-DeleteContainer-3", err)
			return err
		}
	}

	err = StopContainer(containerID)
	if err != nil {
		return err
	}

	err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})

	if err != nil {
		return err
	}

	return nil
}
