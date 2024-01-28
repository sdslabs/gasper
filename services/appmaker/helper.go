package appmaker

import (
	"fmt"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
	"os"
	"path/filepath"
)

var path, _ = os.Getwd()

// storageCleanup removes the application's local storage directory
func storageCleanup(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		utils.LogError("AppMaker-Helper-1", err)
	}
	return err
}

// containerCleanup removes the application's container
func containerCleanup(appName string) error {
	err := docker.DeleteContainer(appName)
	if err != nil {
		utils.LogError("AppMaker-Helper-2", err)
	}
	return err
}

// diskCleanup cleans the specified application's container and local storage
func diskCleanup(appName string) {
	appDir := filepath.Join(path, fmt.Sprintf("storage/%s", appName))
	storeCleanupChan := make(chan error)
	go func() {
		storeCleanupChan <- storageCleanup(appDir)
	}()
	containerCleanup(appName)
	<-storeCleanupChan
}

// stateCleanup removes the application's data from MongoDB and Redis
func stateCleanup(appName string) {
	_, err := mongo.DeleteInstance(types.M{
		mongo.NameKey:         appName,
		mongo.InstanceTypeKey: mongo.AppInstance,
	})
	if err != nil {
		utils.LogError("AppMaker-Helper-3", err)
	}
	if err := redis.RemoveApp(appName); err != nil {
		utils.LogError("AppMaker-Helper-4", err)
	}
}

func fetchAllApplicationNames() []string {
	apps := mongo.FetchDocs(mongo.InstanceCollection, types.M{
		mongo.InstanceTypeKey: mongo.AppInstance,
	})
	var appNames []string
	for _, app := range apps {
		appNames = append(appNames, app[mongo.NameKey].(string))
	}
	return appNames
}