package docker

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// CopyToContainer copies the file from source path to the destination path inside the container
// Reader must be a tar archive
func CopyToContainer(ctx context.Context, cli *client.Client, containerID, destination string, reader io.Reader) error {
	config := types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}
	err := cli.CopyToContainer(ctx, containerID, destination, reader, config)
	if err != nil {
		return err
	}
	return nil
}
