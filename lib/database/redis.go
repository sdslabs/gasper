package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdslabs/gasper/lib/utils"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/types"
)

// CreateRedisDB  creates a RedisDB container
func CreateRedisDB(db types.Database) error {
	port, err := utils.GetFreePort()
	if err != nil {
		return fmt.Errorf("Error while getting free port for container : %s", err)
	}

	storedir := filepath.Join(storepath, "redis-storage", db.GetName())

	if err := os.MkdirAll(storedir, 0755); err != nil {
		return fmt.Errorf("Error while creating the directory : %s", err)
	}

	containerID, err := docker.CreateDatabaseContainer(types.DatabaseContainer{
		Image:         configs.ImageConfig.Redis,
		ContainerPort: port,
		DatabasePort:  6379,
		WorkDir:       "/data/",
		StoreDir:      filepath.Join(storepath, "redis-storage", db.GetName()),
		Name:          db.GetName(),
		Cmd:           []string{"redis-server", "--logfile", "/data/redis-server.log", "--requirepass", db.GetPassword()},
	})

	if err != nil {
		return types.NewResErr(500, "container not created", err)
	}

	if err := docker.StartContainer(containerID); err != nil {
		return types.NewResErr(500, "container not started", err)
	}

	db.SetContainerPort(port)
	return nil
}

// DeleteRedisDB deletes RedisDB container
func DeleteRedisDB(databaseName string) error {
	if err := docker.DeleteContainer(databaseName); err != nil {
		return types.NewResErr(500, "container not deleted", err)
	}

	storedir := filepath.Join(storepath, "redis-storage", databaseName)

	if err := os.RemoveAll(storedir); err != nil {
		return fmt.Errorf("Error while deleting the database directory : %s", err)
	}
	return nil
}

// GetLogs returns logs of a RedisDB container
func GetLogs(databaseName string) (string, error) {
	logs, err := docker.ExecProcessWthStream(databaseName, []string{"sh", "-c", "cat /data/redis-server.log"})
	if err != nil {
		return "", types.NewResErr(500, "Command not executed", err)
	}
	return logs, nil
}
