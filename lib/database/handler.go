package database

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/docker/docker/client"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/types"
	"golang.org/x/net/context"
)

// SetupDBInstance sets up containers for database
func SetupDBInstance(dbtype string) (string, types.ResponseError) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", types.NewResErr(500, "cannot setup client", err)
	}

	storepath, _ := os.Getwd()

	var containerID string

	switch dbtype {
	case mongo.MongoDB:
		{
			dockerImage := configs.ImageConfig.Mongodb
			port := strconv.Itoa(configs.ServiceConfig.Mongodb.ContainerPort)
			env := configs.ServiceConfig.Mongodb.Env
			workdir := "/data/db"
			storedir := filepath.Join(storepath, "mongodb-storage")
			containerID, err = docker.CreateMongoDBContainer(
				ctx,
				cli,
				dockerImage,
				port,
				workdir,
				storedir,
				env)
		}
	case mongo.Mysql:
		{
			dockerImage := configs.ImageConfig.Mysql
			port := strconv.Itoa(configs.ServiceConfig.Mysql.ContainerPort)
			env := configs.ServiceConfig.Mysql.Env
			workdir := "/var/lib/mysql"
			storedir := filepath.Join(storepath, "mysql-storage")
			containerID, err = docker.CreateMysqlContainer(
				ctx,
				cli,
				dockerImage,
				port,
				workdir,
				storedir,
				env)
		}
	}
	err = docker.StartContainer(containerID)

	if err != nil {
		return "", types.NewResErr(500, "container not started", err)
	}

	return containerID, nil
}
