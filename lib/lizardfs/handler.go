package lizardfs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/types"
)

var storepath, _ = os.Getwd()

// Maps lizardfs's service name with its appropriate configuration
var lizardfsMap = map[string]*types.LizardfsContainer{
	types.LizardfsMaster: {
		Image:          "katharostech/lizardfs",
		Cmd:            []string{"master"},
		HostPort1:      9421,
		ContainerPort1: 9421,
		Name:           types.LizardfsMaster,
		WorkDir:        "/var/lib/mfs",
		StoreDir:       filepath.Join(storepath, "lizardfs-storage", "master-storage"),
	},
	types.LizardfsMasterShadow: {
		Image:    "katharostech/lizardfs",
		Cmd:      []string{"master"},
		Env:      map[string]interface{}{"MFSMASTER_PERSONALITY": "shadow"},
		Name:     types.LizardfsMasterShadow,
		WorkDir:  "/var/lib/mfs",
		StoreDir: filepath.Join(storepath, "lizardfs-storage", "master-shadow-storage"),
	},
	types.LizardfsMetalogger: {
		Image:    "katharostech/lizardfs",
		Cmd:      []string{"metalogger"},
		Name:     types.LizardfsMetalogger,
		WorkDir:  "/var/lib/mfs",
		StoreDir: filepath.Join(storepath, "lizardfs-storage", "metalogger-storage"),
	},
	types.LizardfsChunkserver: {
		Image:    "katharostech/lizardfs",
		Cmd:      []string{"chunkserver"},
		Env:      map[string]interface{}{"MFSHDD_1": "/mnt/mfshdd"},
		Name:     types.LizardfsChunkserver,
		WorkDir:  "/mnt/mfshdd",
		StoreDir: filepath.Join(storepath, "lizardfs-storage", "chunkserver-storage"),
	},
	"client1": {
		Image: "katharostech/lizardfs",
		Cmd:   []string{"client", "/mnt/mfs"},
		Name:  "client1",
	},
}

// SetupLizardfsInstance sets up containers for database
func SetupLizardfsInstance(serviceType string) (string, types.ResponseError) {
	if lizardfsMap[serviceType] == nil {
		return "", types.NewResErr(500, fmt.Sprintf("Invalid lizardfs type %s provided", serviceType), nil)
	}
	containerID, err := docker.CreateLizardfsContainer(lizardfsMap[serviceType])
	if err != nil {
		return "", types.NewResErr(500, "container not created", err)
	}

	if err := docker.StartContainer(containerID); err != nil {
		return "", types.NewResErr(500, "container not started", err)
	}

	return containerID, nil
}
