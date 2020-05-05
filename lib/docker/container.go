package docker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/sdslabs/gasper/configs"
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
			volume: {},
		},
		Healthcheck: &container.HealthConfig{
			Test:     []string{"CMD-SHELL", fmt.Sprintf("curl --fail --silent http://localhost:%d/ || exit 1", containerCfg.ApplicationPort)},
			Interval: configs.ServiceConfig.Mizu.MetricsInterval * time.Second,
			Timeout:  10 * time.Second,
			Retries:  3,
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
			volume: {},
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
			volume: {},
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

// CreatePostgreSQLContainer function sets up a postgreSQL instance for managing databases
func CreatePostgreSQLContainer(image, postgresqlPort, workdir, storedir string, env types.M) (string, error) {
	ctx := context.Background()
	volume := fmt.Sprintf("%s:%s", storedir, workdir)

	envArr := []string{}
	for key, value := range env {
		envArr = append(envArr, key+"="+fmt.Sprintf("%v", value))
	}

	containerConfig := &container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			"5432/tcp": struct{}{},
		},
		Env: envArr,
		Volumes: map[string]struct{}{
			volume: {},
		},
	}

	hostConfig := &container.HostConfig{
		Binds: []string{
			volume,
		},
		PortBindings: nat.PortMap{
			nat.Port("5432/tcp"): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: postgresqlPort}},
		},
	}

	createdConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, types.PostgreSQL)
	if err != nil {
		return "", err
	}

	return createdConf.ID, nil
}

// CreateRedisContaienr function sets up a redisDB instance for managing databases
func CreateRedisContainer(image, rediskaenport, workdir, storedir string, containername string, env types.M) (string, error) {
	ctx := context.Background()
	volume := fmt.Sprintf("%s:%s", storedir, workdir)

	envArr := []string{}
	for key, value := range env {
		envArr = append(envArr, key+"="+fmt.Sprintf("%v", value))
	}

	containerConfig := &container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			"6379/tcp": struct{}{},
		},
		Env: envArr,
		Volumes: map[string]struct{}{
			volume: {},
		},
	}

	hostConfig := &container.HostConfig{
		Binds: []string{
			volume,
		},
		PortBindings: nat.PortMap{
			nat.Port("6379/tcp"): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(1000)}},
		},
	}

	createdConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, containername)
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
func ListContainers() ([]string, error) {
	ctx := context.Background()
	containers, err := cli.ContainerList(ctx, dockerTypes.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	list := make([]string, 0)

	for _, container := range containers {
		if len(container.Names) > 0 && len(container.Names[0]) > 1 {
			list = append(list, container.Names[0][1:])
		}
	}
	return list, nil
}

// ContainerStats returns container statistics using the containerID
func ContainerStats(containerID string) (*types.Stats, error) {
	ctx := context.Background()
	containerStats, err := cli.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(containerStats.Body)
	if err != nil {
		return nil, err
	}
	containerStatsInterface := &types.Stats{}
	err = json.Unmarshal(body, containerStatsInterface)
	return containerStatsInterface, err
}
