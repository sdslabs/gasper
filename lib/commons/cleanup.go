package commons

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sdslabs/SWS/lib/database"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
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
func MysqlDatabaseCleanup(dbKey string) error {
	dbUser := strings.Split(dbKey, ":")[0]
	dbName := strings.Split(dbKey, ":")[1]
	return database.DeleteMysqlDB(dbName, dbUser)
}

// MongoDatabaseCleanup cleans the database's space in the container
func MongoDatabaseCleanup(dbKey string) error {
	dbUser := strings.Split(dbKey, ":")[0]
	dbName := strings.Split(dbKey, ":")[1]
	return database.DeleteMongoDB(dbName, dbUser)
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
func DatabaseFullCleanup(dbKey, databaseType string) {
	switch databaseType {
	case mongo.Mysql:
		{
			err := MysqlDatabaseCleanup(dbKey)
			if err != nil {
				utils.LogError(err)
			}
		}
	case mongo.MongoDB:
		{
			err := MongoDatabaseCleanup(dbKey)
			if err != nil {
				utils.LogError(err)
			}
		}
	}
}

// DatabaseStateCleanup removes the database's data from MongoDB and Redis
func DatabaseStateCleanup(dbKey string) {
	dbUser := strings.Split(dbKey, ":")[0]
	dbName := strings.Split(dbKey, ":")[1]
	mongo.DeleteInstance(map[string]interface{}{
		"name":         dbName,
		"user":         dbUser,
		"instanceType": mongo.DBInstance,
	})
	redis.RemoveDB(dbKey)
}
