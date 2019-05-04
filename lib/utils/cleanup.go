package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
)

// StorageCleanup removes the application's local storage directory
func StorageCleanup(path string) error {
	err := os.RemoveAll(path)

	if err != nil {
		return err
	}

	return nil
}

// ContainerCleanup removes the application's container
func ContainerCleanup(appName string) error {
	return docker.DeleteContainer(appName)
}

// FullCleanup cleans the specified application's container and local storage
func FullCleanup(appName string) {
	var (
		path, _ = os.Getwd()
		appDir  = filepath.Join(path, fmt.Sprintf("storage/%s", appName))
	)
	err := StorageCleanup(appDir)
	if err != nil {
		LogError(err)
	}

	err = ContainerCleanup(appName)
	if err != nil {
		LogError(err)
	}
}

// StateCleanup removes the application data from MongoDB and Redis
func StateCleanup(appName string) {
	mongo.DeleteApp(map[string]interface{}{
		"name": appName,
	})

	redis.RemoveApp(appName)
}
