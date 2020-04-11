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

//GetDockerConfig returns the docker conatiner configurations
func GetDockerConfig(image, Port, workdir, storedir string, env types.M, databaseType string) (string, error) {
	switch databaseType {
	case types.MongoDB, types.MongoDBGasper:
		{
			return docker.CreateMongoDBContainer(image, Port, workdir, storedir, env, databaseType)
		}
	case types.MySQL:
		{
			return docker.CreateMysqlContainer(image, Port, workdir, storedir, env, databaseType)
		}
	case types.RedisGasper:
		{
			return docker.CreateRedisContainer(image, Port, workdir, storedir, env, databaseType)
		}
	default:
		{
			return "Invalid database type!", nil
		}
	}
}

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
			containerID, err = GetDockerConfig(dockerImage, port, workdir, storedir, env, databaseType)
		}
	case types.MySQL:
		{
			dockerImage := configs.ImageConfig.Mysql
			port := strconv.Itoa(configs.ServiceConfig.Kaen.MySQL.ContainerPort)
			env := configs.ServiceConfig.Kaen.MySQL.Env
			workdir := "/var/lib/mysql"
			storedir := filepath.Join(storepath, "mysql-storage")
			containerID, err = GetDockerConfig(dockerImage, port, workdir, storedir, env, databaseType)
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
				env,
				databaseType)
		}
	case types.MongoDBGasper:
		{
			dockerImage := configs.ImageConfig.Mongodb
			port := strconv.Itoa(configs.ServiceConfig.Kaze.MongoDB.ContainerPort)
			env := configs.ServiceConfig.Kaze.MongoDB.Env
			workdir := "/data/db"
			storedir := filepath.Join(storepath, "gasper-mongodb-storage")
			containerID, err = GetDockerConfig(dockerImage, port, workdir, storedir, env, databaseType)
		}
	case types.RedisGasper:
		{
			dockerImage := configs.ImageConfig.Redis
			port := strconv.Itoa(configs.ServiceConfig.Kaze.Redis.ContainerPort)
			env := configs.ServiceConfig.Kaze.Redis.Env
			workdir := "/data/db"
			storedir := filepath.Join(storepath, "gasper-redis-storage")
			containerID, err = GetDockerConfig(dockerImage, port, workdir, storedir, env, databaseType)
		}
	default:
		return "", types.NewResErr(500, "invalid database type provided", errors.New("invalid database type provided"))
	}

	if err != nil {
		return "", types.NewResErr(500, "container not created", err)
	}

	err = docker.StartContainer(containerID)

	if err != nil {
		return "", types.NewResErr(500, "container not started", err)
	}

	return containerID, nil
}
