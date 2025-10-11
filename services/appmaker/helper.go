package appmaker

import (
	"os"

	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

var path, _ = os.Getwd()

// storageCleanup removes the application's local storage directory
func storageCleanup(appName string) error {
	err := docker.DeleteStorage(appName)
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
	storeCleanupChan := make(chan error)
	containerCleanup(appName)
	go func() {
		storeCleanupChan <- storageCleanup(appName)
	}()
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

func FetchAllApplicationNames() []string {
	apps := mongo.FetchDocs(mongo.InstanceCollection, types.M{
		mongo.InstanceTypeKey: mongo.AppInstance,
	})
	var appNames []string
	for _, app := range apps {
		appNames = append(appNames, app[mongo.NameKey].(string))
	}
	return appNames
}
