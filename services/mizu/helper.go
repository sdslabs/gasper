package mizu

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// storageCleanup removes the application's local storage directory
func storageCleanup(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		utils.LogError(err)
	}
	return err
}

// containerCleanup removes the application's container
func containerCleanup(appName string) error {
	err := docker.DeleteContainer(appName)
	if err != nil {
		utils.LogError(err)
	}
	return err
}

// diskCleanup cleans the specified application's container and local storage
func diskCleanup(appName string) {
	var (
		path, _ = os.Getwd()
		appDir  = filepath.Join(path, fmt.Sprintf("storage/%s", appName))
	)
	storeCleanupChan := make(chan error)
	containerCleanupChan := make(chan error)
	go func() {
		storeCleanupChan <- storageCleanup(appDir)
	}()
	go func() {
		containerCleanupChan <- containerCleanup(appName)
	}()
	<-storeCleanupChan
	<-containerCleanupChan
}

// stateCleanup removes the application's data from MongoDB and Redis
func stateCleanup(appName string) {
	mongo.DeleteInstance(types.M{
		mongo.NameKey:         appName,
		mongo.InstanceTypeKey: mongo.AppInstance,
	})
	redis.RemoveApp(appName)
}
