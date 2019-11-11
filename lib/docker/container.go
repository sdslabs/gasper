package docker

import (
	"fmt"
	"os"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
	"golang.org/x/net/context"
)

// CreateContainer creates a new container of the given container options, returns id of the container created
func CreateContainer(containerCfg *types.ApplicationContainer) (string, error) {
	ctx := context.Background()
	volume := fmt.Sprintf("%s:%s", containerCfg.StoreDir, containerCfg.WorkDir)

	// convert map to list of strings
	envArr := []string{}
	for key, value := range containerCfg.Env {
		envArr = append(envArr, fmt.Sprintf("%s=%v", key, value))
	}

	containerPortRule := nat.Port(fmt.Sprintf(`%d/tcp`, containerCfg.ApplicationPort))

	containerConfig := &container.Config{
		WorkingDir: containerCfg.WorkDir,
		Image:      containerCfg.Image,
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
			volume,
		},
		DNS: containerCfg.NameServers,
		PortBindings: nat.PortMap{
			nat.Port(containerPortRule): []nat.PortBinding{{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", containerCfg.ContainerPort)}},
		},
		Resources: container.Resources{
			NanoCPUs: containerCfg.CPU,
			Memory:   containerCfg.Memory,
		},
	}

	createdConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, containerCfg.Name)
	if err != nil {
		return "", err
	}
	return createdConf.ID, nil
}

// CreateMysqlContainer function sets up a mysql instance for managing databases
func CreateMysqlContainer(image, mysqlPort, workdir, storedir string, env types.M) (string, error) {
	ctx := context.Background()
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
func CreateMongoDBContainer(image, mongodbPort, workdir, storedir string, env types.M) (string, error) {
	ctx := context.Background()
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
	return cli.ContainerStart(ctx, containerID, dockerTypes.ContainerStartOptions{})
}

// StopContainer stops the container corresponding to given containerID
func StopContainer(containerID string) error {
	ctx := context.Background()
	return cli.ContainerStop(ctx, containerID, nil)
}

// ListContainers lists all containers
func ListContainers() []string {
	ctx := context.Background()
	containers, err := cli.ContainerList(ctx, dockerTypes.ContainerListOptions{All: true})
	if err != nil {
		utils.LogError(err)
		os.Exit(1)
	}

	list := make([]string, 1)

	for _, container := range containers {
		if len(container.Names) > 0 {
			list = append(list, container.Names[0])
		}
	}
	return list
}
