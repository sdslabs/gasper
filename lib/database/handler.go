package database

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/configs"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
	"golang.org/x/net/context"
)

func SetupDBInstance(dbtype string) (string, types.ResponseError) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", types.NewResErr(500, "cannot setup client", err)
	}

	dockerImage := configs.ServiceConfig[dbtype].(map[string]interface{})["image"].(string)
	port := configs.ServiceConfig[dbtype].(map[string]interface{})["container_port"].(string)
	env := configs.ServiceConfig[dbtype].(map[string]interface{})["env"].(map[string]interface{})
	storepath, _ := os.Getwd()

	var containerID string

	switch dbtype {
	case mongo.MongoDB:
		{
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
