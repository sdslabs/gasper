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
var databaseMap = map[string]*types.DatabaseContainer{
	types.MongoDB: {
		Image:         configs.ImageConfig.Mongodb,
		ContainerPort: configs.ServiceConfig.Kaen.MongoDB.ContainerPort,
		DatabasePort:  27017,
		Env:           configs.ServiceConfig.Kaen.MongoDB.Env,
		WorkDir:       "/data/db",
		StoreDir:      filepath.Join(storepath, "mongodb-storage"),
		Name:          types.MongoDB,
	},
	types.MongoDBGasper: {
		Image:         configs.ImageConfig.Mongodb,
		ContainerPort: configs.ServiceConfig.Kaze.MongoDB.ContainerPort,
		DatabasePort:  27017,
		Env:           configs.ServiceConfig.Kaze.MongoDB.Env,
		WorkDir:       "/data/db",
		StoreDir:      filepath.Join(storepath, "gasper-mongodb-storage"),
		Name:          types.MongoDBGasper,
	},
	types.MySQL: {
		Image:         configs.ImageConfig.Mysql,
		ContainerPort: configs.ServiceConfig.Kaen.MySQL.ContainerPort,
		DatabasePort:  3306,
		Env:           configs.ServiceConfig.Kaen.MySQL.Env,
		WorkDir:       "/app",
		StoreDir:      filepath.Join(storepath, "mysql-storage"),
		Name:          types.MySQL,
	},
	types.RedisGasper: {
		Image:         configs.ImageConfig.Redis,
		ContainerPort: configs.ServiceConfig.Kaze.Redis.ContainerPort,
		DatabasePort:  6379,
		Env:           configs.ServiceConfig.Kaze.Redis.Env,
		WorkDir:       "/data/",
		StoreDir:      filepath.Join(storepath, "gasper-redis-storage"),
		Name:          types.RedisGasper,
		Cmd:           []string{"redis-server", "--requirepass", configs.ServiceConfig.Kaze.Redis.Password},
	},
	types.PostgreSQL: {
		Image:         configs.ImageConfig.Postgresql,
		ContainerPort: configs.ServiceConfig.Kaen.PostgreSQL.ContainerPort,
		DatabasePort:  5432,
		Env:           configs.ServiceConfig.Kaen.PostgreSQL.Env,
		WorkDir:       "/var/lib/postgresql/data",
		StoreDir:      filepath.Join(storepath, "postgresql-storage"),
		Name:          types.PostgreSQL,
	},
}

// SetupDBInstance sets up containers for database
func SetupDBInstance(databaseType string) (string, types.ResponseError) {
	if databaseMap[databaseType] == nil {
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
