package database

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/types"
)

type containerHandler struct {
	dockerImage      string
	port             string
	env              types.M
	workdir          string
	storedir         string
	containerID, err func(string, string, string, string, types.M, string) (string, error)
}


var storepath, _ = os.Getwd()

var databaseMap = map[string]*containerHandler{
	types.MongoDB: {
		dockerImage: configs.ImageConfig.Mongodb,
		port:        strconv.Itoa(configs.ServiceConfig.Kaen.MongoDB.ContainerPort),
		env:         configs.ServiceConfig.Kaen.MongoDB.Env,
		workdir:     "/data/db",
		storedir:    filepath.Join(storepath, "mongodb-storage"),
		containerID: docker.CreateMongoDBContainer,
	},
	types.MongoDBGasper: {
		dockerImage: configs.ImageConfig.Mongodb,
		port:        strconv.Itoa(configs.ServiceConfig.Kaze.MongoDB.ContainerPort),
		env:         configs.ServiceConfig.Kaze.MongoDB.Env,
		workdir:     "/data/db",
		storedir:    filepath.Join(storepath, "gasper-mongodb-storage"),
		containerID: docker.CreateMongoDBContainer,
	},
	types.MySQL: {
		dockerImage: configs.ImageConfig.Mysql,
		port:        strconv.Itoa(configs.ServiceConfig.Kaen.MySQL.ContainerPort),
		env:         configs.ServiceConfig.Kaen.MySQL.Env,
		workdir:     "/var/lib/mysql",
		storedir:    filepath.Join(storepath, "mysql-storage"),
		containerID: docker.CreateMysqlContainer,
	},
	types.RedisGasper: {
		dockerImage: configs.ImageConfig.Redis,
		port:        strconv.Itoa(configs.ServiceConfig.Kaze.Redis.ContainerPort),
		env:         configs.ServiceConfig.Kaze.Redis.Env,
		workdir:     "/data/",
		storedir:    filepath.Join(storepath, "gasper-redis-storage"),
		containerID: docker.CreateRedisContainer,
	},
}

// SetupDBInstance sets up containers for database
func SetupDBInstance(databaseType string) (string, types.ResponseError) {

	var containerID string
	var err error

	containerID, err = databaseMap[databaseType].containerID(databaseMap[databaseType].dockerImage,
		databaseMap[databaseType].port,
		databaseMap[databaseType].workdir,
		databaseMap[databaseType].storedir,
		databaseMap[databaseType].env,
		databaseType,
	)

	if err != nil {
		return "", types.NewResErr(500, "container not created", err)
	}

	err = docker.StartContainer(containerID)

	if err != nil {
		return "", types.NewResErr(500, "container not started", err)
	}

	return containerID, nil
}
