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
	"github.com/sdslabs/gasper/lib/utils"
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
			Interval: configs.ServiceConfig.AppMaker.MetricsInterval * time.Second,
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

// CreateDatabaseContainer function creates a new container of the given container options, returns id of the container created
func CreateDatabaseContainer(containerCfg *types.DatabaseContainer) (string, error) {
	ctx := context.Background()
	volume := fmt.Sprintf("%s:%s", containerCfg.StoreDir, containerCfg.WorkDir)

	volumes, err := ListVolumes()
	if err != nil {
		return "", err
	}
	if !utils.Contains(volumes, containerCfg.StoreDir) {
		utils.LogInfo("No %s volume found in host. Creating the volume.", containerCfg.StoreDir)
		volumename, err := CreateVolume(containerCfg.StoreDir, "kadimasolutions/lizardfs-volume-plugin")
		if err != nil {
			utils.Log(fmt.Sprintf("There was a problem creating %s volume.", containerCfg.StoreDir), utils.ErrorTAG)
			utils.LogError(err)
		} else {
			utils.LogInfo("%s volume has been deployed with name:\t%s \n", containerCfg.StoreDir, volumename)
		}
	}

	envArr := []string{}
	for key, value := range containerCfg.Env {
		envArr = append(envArr, fmt.Sprintf("%s=%v", key, value))
	}

	containerPortRule := nat.Port(fmt.Sprintf(`%d/tcp`, containerCfg.DatabasePort))

	containerConfig := &container.Config{
		Image: containerCfg.Image,
		ExposedPorts: nat.PortSet{
			containerPortRule: struct{}{},
		},
		Env: envArr,
		Volumes: map[string]struct{}{
			volume: {},
		},
	}

	if containerCfg.HasCustomCMD() {
		containerConfig.Cmd = containerCfg.Cmd
	}

	hostConfig := &container.HostConfig{
		Binds: []string{
			volume,
		},
		PortBindings: nat.PortMap{
			nat.Port(containerPortRule): []nat.PortBinding{{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", containerCfg.ContainerPort)}},
		},
	}

	createdConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, containerCfg.Name)
	if err != nil {
		return "", err
	}
	return createdConf.ID, nil
}

// CreateLizardfsContainer creates a new container of the given container options, returns id of the container created
func CreateLizardfsContainer(containerCfg *types.LizardfsContainer) (string, error) {
	ctx := context.Background()
	volume := fmt.Sprintf("%s:%s", containerCfg.StoreDir, containerCfg.WorkDir)

	// convert map to list of strings
	envArr := []string{}
	for key, value := range containerCfg.Env {
		envArr = append(envArr, fmt.Sprintf("%s=%v", key, value))
	}

	containerConfig := &container.Config{
		Image: containerCfg.Image,
		Env:   envArr,
	}

	containerConfig.Cmd = containerCfg.Cmd

	hostConfig := &container.HostConfig{}

	if containerCfg.Name != "client1" {
		hostConfig.Binds = []string{
			volume,
		}
		containerConfig.Volumes = map[string]struct{}{
			volume: {},
		}
	}

	if containerCfg.Name == "mfsmaster" {
		containerPortRule1 := nat.Port(fmt.Sprintf(`%d/tcp`, containerCfg.HostPort1))
		containerConfig.ExposedPorts = nat.PortSet{
			containerPortRule1: struct{}{},
		}
		hostConfig.PortBindings = nat.PortMap{
			nat.Port(containerPortRule1): []nat.PortBinding{{
				HostIP:   "",
				HostPort: fmt.Sprintf("%d", containerCfg.ContainerPort1)}},
		}
	}

	if containerCfg.Name == "client1" {
		hostConfig.CapAdd = []string{"SYS_ADMIN"}
		hostConfig.Devices = []container.DeviceMapping{{PathOnHost: "/dev/fuse", PathInContainer: "/dev/fuse", CgroupPermissions: "rwm"}}
		hostConfig.SecurityOpt = []string{"apparmor:unconfined"}
	}

	if containerCfg.Name != "mfsmaster" {
		hostConfig.Links = []string{
			fmt.Sprintf("%s:mfsmaster", "mfsmaster"),
		}
	}

	createdConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, containerCfg.Name)
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
