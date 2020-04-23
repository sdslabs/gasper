package docker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/types"
	"golang.org/x/net/context"
)

// CreateApplicationContainer creates a new container of the given container options, returns id of the container created
func CreateApplicationContainer(containerCfg *types.ApplicationContainer) (string, error) {
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

// createDatabaseContainerConfig function returns the config variables associated with creation of a container.
func createDatabaseContainerConfig(dockerImage, hostPort, workdir, storedir string,
	containerPort nat.Port, env types.M) (*container.Config, *container.HostConfig) {

	volume := fmt.Sprintf("%s:%s", storedir, workdir)

	envArr := []string{}
	for key, value := range env {
		envArr = append(envArr, fmt.Sprintf("%s=%v", key, value))
	}

	containerConfig := &container.Config{
		Image: dockerImage,
		ExposedPorts: nat.PortSet{
			containerPort: struct{}{},
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
			nat.Port(containerPort): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: hostPort}},
		},
	}

	return containerConfig, hostConfig
}

// createDatabaseContainer creates a container for database services
func createDatabaseContainer(image, mysqlPort, workdir, storedir, databaseType string,
	containerPort nat.Port, env types.M) (string, error) {

	ctx := context.Background()
	containerConfig, hostConfig := createDatabaseContainerConfig(image, mysqlPort, workdir, storedir, containerPort, env)
	createdConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, databaseType)
	if err != nil {
		return "", err
	}
	return createdConf.ID, nil
}

// CreateMySQLContainer function sets up a mysql instance for managing databases
func CreateMySQLContainer(image, mysqlPort, workdir, storedir, databaseType string, env types.M) (string, error) {
	return createDatabaseContainer(image, mysqlPort, workdir, storedir, databaseType, "3306/tcp", env)
}

// CreateMongoDBContainer function sets up a mongoDB instance for managing databases
func CreateMongoDBContainer(image, mongodbPort, workdir, storedir, databaseType string, env types.M) (string, error) {
	return createDatabaseContainer(image, mongodbPort, workdir, storedir, databaseType, "27017/tcp", env)
}

// CreatePostgreSQLContainer function sets up a postgreSQL instance for managing databases
func CreatePostgreSQLContainer(image, postgresqlPort, workdir, storedir, databaseType string, env types.M) (string, error) {
	return createDatabaseContainer(image, postgresqlPort, workdir, storedir, databaseType, "5432/tcp", env)
}

// CreateRedisContainer function sets up a redis instance for gasper
func CreateRedisContainer(image, redisPort, workdir, storedir, databaseType string, env types.M) (string, error) {
	ctx := context.Background()
	containerConfig, hostConfig := createDatabaseContainerConfig(image, redisPort, workdir, storedir, "6379/tcp", env)
	containerConfig.Cmd = []string{"redis-server", "--requirepass", configs.ServiceConfig.Kaze.Redis.Password}
	createdConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, databaseType)
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
