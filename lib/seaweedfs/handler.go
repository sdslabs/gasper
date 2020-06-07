package seaweedfs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/types"
)

var storepath, _ = os.Getwd()

// Maps seaweedfs's service name with its appropriate configuration
var seaweedfsMap = map[string]*types.SeaweedfsContainer{
	types.SeaweedMaster: {
		Image:          "chrislusf/seaweedfs",
		Cmd:            []string{"master", "-ip=master"},
		HostPort1:      9333,
		ContainerPort1: 9333,
		HostPort2:      19333,
		ContainerPort2: 1933,
		WorkDir:        "/data",
		StoreDir:       filepath.Join(storepath, "seaweed", "seaweed-master-storage"),
		Name:           types.SeaweedMaster,
	},
	types.SeaweedVolume: {
		Image:          "chrislusf/seaweedfs",
		Cmd:            []string{"volume", "-mserver=master:9333", "-port=8080"},
		HostPort1:      8080,
		ContainerPort1: 8080,
		HostPort2:      18080,
		ContainerPort2: 18080,
		WorkDir:        "/data",
		StoreDir:       filepath.Join(storepath, "seaweed", "seaweed-volume-storage"),
		Name:           types.SeaweedVolume,
	},
	types.SeaweedFiler: {
		Image:          "chrislusf/seaweedfs",
		Cmd:            []string{"filer", "-master=master:9333"},
		HostPort1:      8888,
		ContainerPort1: 8888,
		HostPort2:      18888,
		ContainerPort2: 18888,
		WorkDir:        "/data",
		StoreDir:       filepath.Join(storepath, "seaweed", "seaweed-filer-storage"),
		Name:           types.SeaweedFiler,
	},
	types.SeaweedCronjob: {
		Image:          "chrislusf/seaweedfs",
		Cmd:            []string{"cronjob"},
		HostPort1:      8889,
		ContainerPort1: 8889,
		HostPort2:      18889,
		ContainerPort2: 18889,
		WorkDir:        "/data",
		StoreDir:       filepath.Join(storepath, "seaweed", "seaweed-cronjob-storage"),
		Env:            map[string]interface{}{"CRON_SCHEDULE": "*/2 * * * * *", "WEED_MASTER": "master:9333"},
		Name:           types.SeaweedCronjob,
	},
	types.SeaweedS3: {
		Image:          "chrislusf/seaweedfs",
		Cmd:            []string{"s3", "-filer=filer:8888"},
		HostPort1:      8333,
		ContainerPort1: 8333,
		HostPort2:      18898,
		ContainerPort2: 18898,
		WorkDir:        "/data",
		StoreDir:       filepath.Join(storepath, "seaweed", "seaweed-s3-storage"),
		Name:           types.SeaweedS3,
	},
}

// SetupSeaweedfsInstance sets up containers for database
func SetupSeaweedfsInstance(seaweedType string) (string, types.ResponseError) {
	if seaweedfsMap[seaweedType] == nil {
		return "", types.NewResErr(500, fmt.Sprintf("Invalid seaweedfs type %s provided", seaweedType), nil)
	}

	containerID, err := docker.CreateSeaweedContainer(seaweedfsMap[seaweedType])
	if err != nil {
		return "", types.NewResErr(500, "container not created", err)
	}

	if err := docker.StartContainer(containerID); err != nil {
		return "", types.NewResErr(500, "container not started", err)
	}

	return containerID, nil
}
