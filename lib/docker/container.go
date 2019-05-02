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
func CreateContainer(ctx context.Context, cli *client.Client, image, httpPort, sshPort, workdir, storedir, name string, env map[string]interface{}) (string, error) {
	volume := fmt.Sprintf("%s:%s", storedir, workdir)

	// convert map to list of strings
	envArr := []string{}
	for key, value := range env {
		envArr = append(envArr, key+"="+fmt.Sprintf("%v", value))
	}

	containerConfig := &container.Config{
		WorkingDir: workdir,
		Image:      image,
		ExposedPorts: nat.PortSet{
			"80/tcp": struct{}{},
			"22/tcp": struct{}{},
		},
		Env: envArr,
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

// StopContainer stops the container corresponding to given containerID
func StopContainer(ctx context.Context, cli *client.Client, containerID string) error {
	return cli.ContainerStop(ctx, containerID, nil)
}
