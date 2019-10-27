package docker

import (
	"io"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// CopyToContainer copies the file from source path to the destination path inside the container
// Reader must be a tar archive
func CopyToContainer(containerID, destination string, reader io.Reader) error {
	ctx := context.Background()
	config := types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}
	err := cli.CopyToContainer(ctx, containerID, destination, reader, config)
	if err != nil {
		return err
	}
	return nil
}
