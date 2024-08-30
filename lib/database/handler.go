package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/types"
)

var storepath, _ = os.Getwd()

// Maps database name with its appropriate configuration
var databaseMap = map[string]types.DatabaseContainer{
	types.MongoDB: {
		Image:         configs.ImageConfig.Mongodb,
		ContainerPort: configs.ServiceConfig.DbMaker.MongoDB.ContainerPort,
		DatabasePort:  27017,
		Env:           configs.ServiceConfig.DbMaker.MongoDB.Env,
		WorkDir:       "/data/db",
		StoreDir:      filepath.Join(storepath, "mongodb-storage"),
		Name:          types.MongoDB,
	},
	types.MongoDBGasper: {
		Image:         configs.ImageConfig.Mongodb,
		ContainerPort: configs.ServiceConfig.Master.MongoDB.ContainerPort,
		DatabasePort:  27017,
		Env:           configs.ServiceConfig.Master.MongoDB.Env,
		WorkDir:       "/data/db",
		StoreDir:      filepath.Join(storepath, "gasper-mongodb-storage"),
		Name:          types.MongoDBGasper,
	},
	types.MySQL: {
		Image:         configs.ImageConfig.Mysql,
		ContainerPort: configs.ServiceConfig.DbMaker.MySQL.ContainerPort,
		DatabasePort:  3306,
		Env:           configs.ServiceConfig.DbMaker.MySQL.Env,
		WorkDir:       "/app",
		StoreDir:      filepath.Join(storepath, "mysql-storage"),
		Name:          types.MySQL,
	},
	types.RedisGasper: {
		Image:         configs.ImageConfig.Redis,
		ContainerPort: configs.ServiceConfig.Master.Redis.ContainerPort,
		DatabasePort:  6379,
		WorkDir:       "/data/",
		StoreDir:      filepath.Join(storepath, "gasper-redis-storage"),
		Name:          types.RedisGasper,
		Cmd:           []string{"redis-server", "--requirepass", configs.ServiceConfig.Master.Redis.Password},
	},
	types.PostgreSQL: {
		Image:         configs.ImageConfig.Postgresql,
		ContainerPort: configs.ServiceConfig.DbMaker.PostgreSQL.ContainerPort,
		DatabasePort:  5432,
		Env:           configs.ServiceConfig.DbMaker.PostgreSQL.Env,
		WorkDir:       "/var/lib/postgresql/data",
		StoreDir:      filepath.Join(storepath, "postgresql-storage"),
		Name:          types.PostgreSQL,
	},
}

// SetupDBInstance sets up containers for database
func SetupDBInstance(databaseType string) (string, types.ResponseError) {
	if _, found := databaseMap[databaseType]; !found {
		return "", types.NewResErr(500, fmt.Sprintf("Invalid database type %s provided", databaseType), nil)
	}

	containerID, err := docker.CreateDatabaseContainer(databaseMap[databaseType])
	if err != nil {
		return "", types.NewResErr(500, "container not created", err)
	}

	if err := docker.StartContainer(containerID); err != nil {
		return "", types.NewResErr(500, "container not started", err)
	}

	return containerID, nil
}

// LogDB logs the database logs (tail 10) after metrics interval i.e. 1 minute
func LogDB(service string) (string, error) {
	var log_location string
	switch service {
	case types.MySQL:
		log_location = "/var/log/mysql/general.log"

	case types.PostgreSQL:
		log_location = "/var/lib/postgresql/data/pg_log/postgresql_log.log"

	case types.MongoDB:
		log_location = "/var/log/mongodb/mongodb.log"
	}
	log_string, err := docker.ExecProcessWthStream(service, []string{"sh", "-c", fmt.Sprintf("cat %s", log_location)})
	if err != nil {
		return "", err
	}
	return log_string, nil
}
