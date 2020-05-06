package commons

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdslabs/gasper/lib/database"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// StorageCleanup removes the application's local storage directory
func StorageCleanup(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		utils.LogError(err)
	}
	return err
}

// ContainerCleanup removes the application's container
func ContainerCleanup(appName string) error {
	err := docker.DeleteContainer(appName)
	if err != nil {
		utils.LogError(err)
	}
	return err
}

// MysqlDatabaseCleanup cleans the database's space in the container
func MysqlDatabaseCleanup(db string) error {
	return database.DeleteMysqlDB(db)
}

// MongoDatabaseCleanup cleans the database's space in the container
func MongoDatabaseCleanup(db string) error {
	return database.DeleteMongoDB(db)
}

// AppFullCleanup cleans the specified application's container and local storage
func AppFullCleanup(instanceName string) {
	var (
		path, _ = os.Getwd()
		appDir  = filepath.Join(path, fmt.Sprintf("storage/%s", instanceName))
	)
	storeCleanupChan := make(chan error)
	containerCleanupChan := make(chan error)
	go func() {
		storeCleanupChan <- StorageCleanup(appDir)
	}()
	go func() {
		containerCleanupChan <- ContainerCleanup(instanceName)
	}()
	<-storeCleanupChan
	<-containerCleanupChan
}

// AppStateCleanup removes the application's data from MongoDB and Redis
func AppStateCleanup(instanceName string) {
	mongo.DeleteInstance(types.M{
		mongo.NameKey:         instanceName,
		mongo.InstanceTypeKey: mongo.AppInstance,
	})
	redis.RemoveApp(instanceName)
}
