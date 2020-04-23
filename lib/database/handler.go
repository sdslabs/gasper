package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/types"
)

var storepath, _ = os.Getwd()

type containerHandler struct {
	dockerImage           string
	port                  string
	env                   types.M
	workdir               string
	storedir              string
	databaseType          string
	containerSpawner, err func(string, string, string, string, string, types.M) (string, error)
}

func (ch *containerHandler) spawn() (string, error) {
	return ch.containerSpawner(
		ch.dockerImage,
		ch.port,
		ch.workdir,
		ch.storedir,
		ch.databaseType,
		ch.env,
	)
}

var databaseMap = map[string]*containerHandler{
	types.MongoDB: {
		dockerImage:      configs.ImageConfig.Mongodb,
		port:             strconv.Itoa(configs.ServiceConfig.Kaen.MongoDB.ContainerPort),
		env:              configs.ServiceConfig.Kaen.MongoDB.Env,
		workdir:          "/data/db",
		storedir:         filepath.Join(storepath, "mongodb-storage"),
		databaseType:     types.MongoDB,
		containerSpawner: docker.CreateMongoDBContainer,
	},
	types.MongoDBGasper: {
		dockerImage:      configs.ImageConfig.Mongodb,
		port:             strconv.Itoa(configs.ServiceConfig.Kaze.MongoDB.ContainerPort),
		env:              configs.ServiceConfig.Kaze.MongoDB.Env,
		workdir:          "/data/db",
		storedir:         filepath.Join(storepath, "gasper-mongodb-storage"),
		databaseType:     types.MongoDBGasper,
		containerSpawner: docker.CreateMongoDBContainer,
	},
	types.MySQL: {
		dockerImage:      configs.ImageConfig.Mysql,
		port:             strconv.Itoa(configs.ServiceConfig.Kaen.MySQL.ContainerPort),
		env:              configs.ServiceConfig.Kaen.MySQL.Env,
		workdir:          "/var/lib/mysql",
		storedir:         filepath.Join(storepath, "mysql-storage"),
		databaseType:     types.MySQL,
		containerSpawner: docker.CreateMySQLContainer,
	},
	types.RedisGasper: {
		dockerImage:      configs.ImageConfig.Redis,
		port:             strconv.Itoa(configs.ServiceConfig.Kaze.Redis.ContainerPort),
		env:              configs.ServiceConfig.Kaze.Redis.Env,
		workdir:          "/data/",
		storedir:         filepath.Join(storepath, "gasper-redis-storage"),
		databaseType:     types.RedisGasper,
		containerSpawner: docker.CreateRedisContainer,
	},
	types.PostgreSQL: {
		dockerImage:      configs.ImageConfig.Postgresql,
		port:             strconv.Itoa(configs.ServiceConfig.Kaen.PostgreSQL.ContainerPort),
		env:              configs.ServiceConfig.Kaen.PostgreSQL.Env,
		workdir:          "/var/lib/postgresql/data",
		storedir:         filepath.Join(storepath, "postgresql-storage"),
		databaseType:     types.PostgreSQL,
		containerSpawner: docker.CreatePostgreSQLContainer,
	},
}

// SetupDBInstance sets up containers for database
func SetupDBInstance(databaseType string) (string, types.ResponseError) {
	if databaseMap[databaseType] == nil {
		return "", types.NewResErr(500, fmt.Sprintf("Invalid database type %s provided", databaseType), nil)
	}

	containerID, err := databaseMap[databaseType].spawn()
	if err != nil {
		return "", types.NewResErr(500, "container not created", err)
	}

	if err := docker.StartContainer(containerID); err != nil {
		return "", types.NewResErr(500, "container not started", err)
	}

	return containerID, nil
}
