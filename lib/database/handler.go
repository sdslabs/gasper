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
			port := strconv.Itoa(configs.ServiceConfig.Mongodb.ContainerPort)
			env := configs.ServiceConfig.Mongodb.Env
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
			port := strconv.Itoa(configs.ServiceConfig.Mysql.ContainerPort)
			env := configs.ServiceConfig.Mysql.Env
			workdir := "/var/lib/mysql"
			storedir := filepath.Join(storepath, "mysql-storage")
			containerID, err = docker.CreateMysqlContainer(
				dockerImage,
				port,
				workdir,
				storedir,
				env)
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
