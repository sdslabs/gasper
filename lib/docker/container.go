package docker

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
)

// CreateContainer creates a new container of the given container options, returns id of the container created
func CreateContainer(ctx context.Context, cli *client.Client, image, httpPort, sshPort, workdir, storedir, name string) (string, error) {
	volume := fmt.Sprintf("%s:%s", storedir, workdir)
	containerConfig := &container.Config{
		WorkingDir: workdir,
		Image:      image,
		ExposedPorts: nat.PortSet{
			"80/tcp": struct{}{},
			"22/tcp": struct{}{},
		},
		Volumes: map[string]struct{}{
			volume: struct{}{},
		},
	}
	hostConfig := &container.HostConfig{
		Binds: []string{
			"/var/run/docker.sock:/var/run/docker.sock",
			volume,
		},
		PortBindings: nat.PortMap{
			nat.Port("80/tcp"): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: httpPort}},
			nat.Port("22/tcp"): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: sshPort}},
		},
	}
	createdConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, name)
	if err != nil {
		return "", err
	}
	return createdConf.ID, nil
}

// StartContainer starts the container corresponding to given containerID
func StartContainer(ctx context.Context, cli *client.Client, containerID string) error {
	return cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}
