package docker

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
)

// CreateContainer spawns a new container of the provided docker image
func CreateContainer(image string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// Map 0.0.0.0:7000 -> 80/tcp
	config := &container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			"80/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		Binds: []string{
			"/var/run/docker.sock:/var/run/docker.sock",
		},
		PortBindings: nat.PortMap{
			nat.Port("80/tcp"): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "7000"}},
		},
	}
	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, "static")

	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

// AddFileToContainer copies the file from source path to the destination path inside the container
// Reader must be a tar archive
func AddFileToContainer(containerID, destination string, reader io.Reader) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	config := types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}
	err = cli.CopyToContainer(ctx, containerID, destination, reader, config)
	if err != nil {
		return err
	}
	return nil
}
