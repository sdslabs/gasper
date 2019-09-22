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

<<<<<<< HEAD
// DatabaseCleanup cleans the database's space in the container
func DatabaseCleanup(dbKey string) error {
	dbUser := strings.Split(dbKey, ":")[0]
	dbName := strings.Split(dbKey, ":")[1]
	return database.DeleteDB(dbName, dbUser)
=======
// MysqlDatabaseCleanup cleans the database's space in the container
func MysqlDatabaseCleanup(dbKey string) error {
	dbName := strings.Split(dbKey, ":")[0]
	dbUser := strings.Split(dbKey, ":")[1]
	return database.DeleteMysqlDB(dbName, dbUser)
}

// MongoDatabaseCleanup cleans the database's space in the container
func MongoDatabaseCleanup(dbKey string) error {
	dbName := strings.Split(dbKey, ":")[0]
	dbUser := strings.Split(dbKey, ":")[1]
	dbPass := strings.Split(dbKey, ":")[2]
	return database.DeleteMongoDB(dbName, dbUser, dbPass)
>>>>>>> checked working
}

// FullCleanup cleans the specified application's container and local storage
func FullCleanup(instanceName, instanceType string) {
	switch instanceType {
	case mongo.AppInstance:
		{
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
<<<<<<< HEAD
	case mongo.DBInstance:
=======
	case "mysqldb":
		{
			err := MysqlDatabaseCleanup(instanceName)
			if err != nil {
				utils.LogError(err)
			}
		}
	case "mongoDb":
>>>>>>> checked working
		{
			err := MongoDatabaseCleanup(instanceName)
			if err != nil {
				utils.LogError(err)
			}
		}
	}
}

// StateCleanup removes the application's/database's data from MongoDB and Redis
func StateCleanup(instanceName, instanceType string) {
	mongo.DeleteInstance(map[string]interface{}{
		"name":         instanceName,
		"instanceType": instanceType,
	})

	switch instanceType {
	case mongo.AppInstance:
		redis.RemoveApp(instanceName)
<<<<<<< HEAD
	case mongo.DBInstance:
=======
	case "mysqldb":
		redis.RemoveDB(instanceName)
	case "mongoDb":
>>>>>>> checked working
		redis.RemoveDB(instanceName)
	}
}
