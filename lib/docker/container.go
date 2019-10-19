package docker

import (
	"fmt"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
	"golang.org/x/net/context"
)

// CreateContainer creates a new container of the given container options, returns id of the container created
func CreateContainer(
	ctx context.Context,
	cli *client.Client,
	image, httpPort, containerPort, workdir, storedir, name string,
	resources container.Resources,
	env types.M) (string, error) {

	volume := fmt.Sprintf("%s:%s", storedir, workdir)

	// convert map to list of strings
	envArr := []string{}
	for key, value := range env {
		envArr = append(envArr, key+"="+fmt.Sprintf("%v", value))
	}

	var containerPortRule nat.Port
	containerPortRule = nat.Port(fmt.Sprintf(`%s/tcp`, containerPort))

	containerConfig := &container.Config{
		WorkingDir: workdir,
		Image:      image,
		ExposedPorts: nat.PortSet{
			containerPortRule: struct{}{},
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
		DNS: configs.GasperConfig.DNSServers,
		PortBindings: nat.PortMap{
			nat.Port(containerPortRule): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: httpPort}},
		},
		Resources: resources,
	}

	createdConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, name)
	if err != nil {
		return "", err
	}
	return createdConf.ID, nil
}

// CreateMysqlContainer function sets up a mysql instance for managing databases
func CreateMysqlContainer(ctx context.Context, cli *client.Client, image, mysqlPort, workdir, storedir string, env types.M) (string, error) {
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

	createdConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, types.MySQL)

	if err != nil {
		return "", err
	}

	return createdConf.ID, nil
}

// CreateMongoDBContainer function sets up a mongoDB instance for managing databases
func CreateMongoDBContainer(ctx context.Context, cli *client.Client, image, mongodbPort, workdir, storedir string, env types.M) (string, error) {
	volume := fmt.Sprintf("%s:%s", storedir, workdir)

	envArr := []string{}
	for key, value := range env {
		envArr = append(envArr, key+"="+fmt.Sprintf("%v", value))
	}

	containerConfig := &container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			"27017/tcp": struct{}{},
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
			nat.Port("27017/tcp"): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: mongodbPort}},
		},
	}

	createdConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, types.MongoDB)
	if err != nil {
		return "", err
	}

	return createdConf.ID, nil
}

// StartContainer starts the container corresponding to given containerID
func StartContainer(containerID string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		utils.LogError(err)
		return err
	}
	return cli.ContainerStart(ctx, containerID, dockerTypes.ContainerStartOptions{})
}

// StopContainer stops the container corresponding to given containerID
func StopContainer(containerID string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		utils.LogError(err)
		return err
	}
	return cli.ContainerStop(ctx, containerID, nil)
}

// ListContainers lists all containers
func ListContainers() []string {
	cli, err := client.NewEnvClient()
	if err != nil {
		utils.LogError(err)
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), dockerTypes.ContainerListOptions{All: true})
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
