package docker

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sdslabs/SWS/lib/utils"
	"golang.org/x/net/context"
)

// CreateContainer creates a new container of the given container options, returns id of the container created
func CreateContainer(ctx context.Context, cli *client.Client, image, httpPort, workdir, storedir, name string, env map[string]interface{}) (string, error) {
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
		},
	}
	createdConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, name)
	if err != nil {
		return "", err
	}
	return createdConf.ID, nil
}

// CreateMysqlContainer function sets up a mysql instance for managing databases
func CreateMysqlContainer(ctx context.Context, cli *client.Client, image, mysqlPort, workdir, storedir string, env map[string]interface{}) (string, error) {
	volume := fmt.Sprintf("%s:%s", storedir, workdir)

	envArr := []string{}
	for key, value := range env {
		envArr = append(envArr, key+"="+fmt.Sprintf("%v", value))
	}

	containerConfig := &container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			"3306/tcp": struct{}{},
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
			nat.Port("3306/tcp"): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: mysqlPort}},
		},
	}

	createdConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, "mysql")
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

func ListContainers() []string {
	cli, err := client.NewEnvClient()
	if err != nil {
		utils.LogError(err)
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		utils.LogError(err)
		panic(err)
	}

	list := make([]string, 1)

	for _, container := range containers {
		if len(container.Names) > 0 {
			list = append(list, container.Names[0])
		}
	}
	return list
}
