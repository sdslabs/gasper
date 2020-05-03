package database

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/types"
)

// SetupDBInstance sets up containers for database
func SetupDBInstance(databaseType string) (string, types.ResponseError) {
	storepath, _ := os.Getwd()

	var containerID string
	var err error

	switch databaseType {
	case types.MongoDB:
		{
			dockerImage := configs.ImageConfig.Mongodb
			port := strconv.Itoa(configs.ServiceConfig.Kaen.MongoDB.ContainerPort)
			env := configs.ServiceConfig.Kaen.MongoDB.Env
			workdir := "/data/db"
			storedir := filepath.Join(storepath, "mongodb-storage")
			containerID, err = docker.CreateMongoDBContainer(
				dockerImage,
				port,
				workdir,
				storedir,
				env)
		}
	case types.MySQL:
		{
			dockerImage := configs.ImageConfig.Mysql
			port := strconv.Itoa(configs.ServiceConfig.Kaen.MySQL.ContainerPort)
			env := configs.ServiceConfig.Kaen.MySQL.Env
			workdir := "/var/lib/mysql"
			storedir := filepath.Join(storepath, "mysql-storage")
			containerID, err = docker.CreateMysqlContainer(
				dockerImage,
				port,
				workdir,
				storedir,
				env)
		}
	case types.PostgreSQL:
		{
			dockerImage := configs.ImageConfig.Postgresql
			port := strconv.Itoa(configs.ServiceConfig.Kaen.PostgreSQL.ContainerPort)
			env := configs.ServiceConfig.Kaen.PostgreSQL.Env
			workdir := "/var/lib/postgresql/data"
			storedir := filepath.Join(storepath, "postgresql-storage")
			containerID, err = docker.CreatePostgreSQLContainer(
				dockerImage,
				port,
				workdir,
				storedir,
				env)
		}
	case types.RedisKaen:
		{
			// dockerImage := configs.ImageConfig.RedisKaen
			// port := strconv.Itoa(configs.ServiceConfig.Kaen.RedisKaen.ContainerPort)
			// env := configs.ServiceConfig.Kaen.RedisKaen.Env
			// workdir := "/var/lib/redis/6379"
			// storedir := filepath.Join(storepath, "redis-storage", "main")
			// containerID, err = docker.CreateRedisContainer(
			// 	dockerImage,
			// 	port,
			// 	workdir,
			// 	storedir,
			// 	"main",
			// 	env)
		}

	default:
		return "", types.NewResErr(500, "invalid database type provided", errors.New("invalid database type provided"))
	}

	if err != nil {
		return "", types.NewResErr(500, "container not created", err)
	}

	if databaseType != types.RedisKaen {
		err = docker.StartContainer(containerID)
	}

	if err != nil {
		return "", types.NewResErr(500, "container not started", err)
	}

	return containerID, nil
}
