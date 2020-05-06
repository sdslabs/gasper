package database

import (
	"os"
	"path/filepath"

	"github.com/sdslabs/gasper/lib/utils"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/types"
)

// CreateRedisDBContainer  creates a RedisDB container
func CreateRedisDBContainer(db types.Database) error {
	storepath, _ := os.Getwd()
	var err error
	port, err := utils.GetFreePort()
	storedir := filepath.Join(storepath, "redis-storage", db.GetName())
	err = os.MkdirAll(storedir, 0755)
	containerID, err := docker.CreateDatabaseContainer(&types.DatabaseContainer{
		Image:         configs.ImageConfig.Redis,
		ContainerPort: port,
		DatabasePort:  6379,
		Env:           configs.ServiceConfig.Kaen.RedisKaen.Env,
		WorkDir:       "/data/",
		StoreDir:      filepath.Join(storepath, "kaen-redis-storage", db.GetName()),
		Name:          db.GetName(),
		Cmd:           []string{"redis-server", "--requirepass", db.GetPassword()},
	})

	err = docker.StartContainer(containerID)
	db.SetContainerPort(port)

	if err != nil {
		return types.NewResErr(500, "container not created", err)
	}

	err = docker.StartContainer(containerID)
	return nil
}

// DeleteRedisDBContainer deletes RedisDB container
func DeleteRedisDBContainer(containerID string) error {
	err := docker.DeleteContainer(containerID)
	storepath, _ := os.Getwd()
	storedir := filepath.Join(storepath, "redis-storage", containerID)
	os.RemoveAll(storedir)
	return err
}
