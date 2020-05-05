package database

import (
	"os"
	"path/filepath"

	"github.com/sdslabs/gasper/lib/utils"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/types"
)

// CreateRedisDBContainer  creates a postgre database
func CreateRedisDBContainer(db types.Database) error {
	storepath, _ := os.Getwd()
	var err error
	port, err := utils.GetFreePort()
	var databaseMap = map[string]*types.DatabaseContainer{
		types.RedisKaen: {
			Image:         configs.ImageConfig.Redis,
			ContainerPort: port,
			DatabasePort:  6379,
			Env:           configs.ServiceConfig.Kaen.RedisKaen.Env,
			WorkDir:       "/data/",
			StoreDir:      filepath.Join(storepath, "kaen-redis-storage", db.GetName()),
			Name:          db.GetName(),
			Cmd:           []string{"redis-server", "--requirepass", db.GetPassword()},
		},
	}

	containerID, err := docker.CreateDatabaseContainer(databaseMap[types.RedisKaen])
	err = docker.StartContainer(containerID)

	db.SetContainerPort(port)

	if err != nil {
		return types.NewResErr(500, "container not created", err)
	}

	err = docker.StartContainer(containerID)

	return nil
}

// DeleteRedisDBContainer delete container
func DeleteRedisDBContainer(containerID string) error {
	err := docker.DeleteContainer(containerID)
	return err

}
