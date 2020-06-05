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
		StoreDir:      "weed-vol", /*filepath.Join(storepath, "mongodb-storage")*/
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

var seaweedfsMap = map[string]*types.SeaweedfsContainer{
	"SeaweedMaster": {
		Image:          "chrislusf/seaweedfs",
		Cmd:            []string{"master", "-ip=master"},
		HostPort1:      9333,
		ContainerPort1: 9333,
		HostPort2:      19333,
		ContainerPort2: 1933,
		WorkDir:        "",
		StoreDir:/*filepath.Join(storepath, "seaweed-master-storage")*/ "weed-voleee",
		Name: "master",
	},
	"SeaweedVolume": {
		Image:          "chrislusf/seaweedfs",
		Cmd:            []string{"volume", "-mserver=master:9333", "-port=8080"},
		HostPort1:      8080,
		ContainerPort1: 8080,
		HostPort2:      18080,
		ContainerPort2: 18080,
		WorkDir:        "",
		StoreDir:       filepath.Join(storepath, "seaweed-volume-storage"),
		Name:           "volume",
	},
	"SeaweedFiler": {
		Image:          "chrislusf/seaweedfs",
		Cmd:            []string{"filer", "-master=master:9333"},
		HostPort1:      8888,
		ContainerPort1: 8888,
		HostPort2:      18888,
		ContainerPort2: 18888,
		WorkDir:        "/data/",
		StoreDir:       filepath.Join(storepath, "seaweed-filer-storage"),
		Name:           "filer",
	},
	"SeaweedCronjob": {
		Image:          "chrislusf/seaweedfs",
		Cmd:            []string{"cronjob"},
		HostPort1:      8889,
		ContainerPort1: 8889,
		HostPort2:      18889,
		ContainerPort2: 18889,
		WorkDir:        "/data/",
		StoreDir:       filepath.Join(storepath, "seaweed-cronjob-storage"),
		Env:            map[string]interface{}{"CRON_SCHEDULE": "*/2 * * * * *", "WEED_MASTER": "master:9333"},
	},
	"SeaweedS3": {
		Image:          "chrislusf/seaweedfs",
		Cmd:            []string{"s3", "-filer=filer:8888"},
		HostPort1:      8333,
		ContainerPort1: 8333,
		HostPort2:      18898,
		ContainerPort2: 18898,
		WorkDir:        "/data/",
		StoreDir:       filepath.Join(storepath, "seaweed-s3-storage"),
	},
	"Seaweed": {
		Image:          "chrislusf/seaweedfs",
		Cmd:            []string{"server", "-filer=true"},
		HostPort1:      9333,
		ContainerPort1: 9333,
		HostPort2:      19333,
		ContainerPort2: 1933,
		WorkDir:        "",
		StoreDir:       filepath.Join(storepath, "seaweed-master-storage"),
		Name:           "seaweed",
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
