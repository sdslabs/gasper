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
	err := StorageCleanup(appDir)
	if err != nil {
		utils.LogError(err)
	}

	err = ContainerCleanup(instanceName)
	if err != nil {
		utils.LogError(err)
	}
}

// AppStateCleanup removes the application's data from MongoDB and Redis
func AppStateCleanup(instanceName string) {
	mongo.DeleteInstance(map[string]interface{}{
		"name":         instanceName,
		"instanceType": mongo.AppInstance,
	})
	redis.RemoveApp(instanceName)
}

// DatabaseFullCleanup deletes the specified database from the container
func DatabaseFullCleanup(db, databaseType string) {
	switch databaseType {
	case mongo.Mysql:
		{
			err := MysqlDatabaseCleanup(db)
			if err != nil {
				utils.LogError(err)
			}
		}
	case mongo.MongoDB:
		{
			err := MongoDatabaseCleanup(db)
			if err != nil {
				utils.LogError(err)
			}
		}
	}
}

// DatabaseStateCleanup removes the database's data from MongoDB and Redis
func DatabaseStateCleanup(db string) {
	mongo.DeleteInstance(map[string]interface{}{
		"name":         db,
		"user":         db,
		"instanceType": mongo.DBInstance,
	})
	redis.RemoveDB(db)
}
