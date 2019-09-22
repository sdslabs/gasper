package database

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
	"golang.org/x/net/context"
)

var dbctx context.Context
var cli *client.Client
var containerID string

// SetupDBInstance sets up the mysql instance for deployment
func SetupDBInstance(dbtype string) (string, types.ResponseError) {
	dbctx = context.Background()
	var err error
	cli, err = client.NewEnvClient()
	if err != nil {
		return "", types.NewResErr(500, "cannot setup client", err)
	}

	dockerImage := configs.ServiceConfig[dbtype].(map[string]interface{})["image"].(string)
	port := configs.ServiceConfig[dbtype].(map[string]interface{})["container_port"].(string)
	env := configs.ServiceConfig[dbtype].(map[string]interface{})["env"].(map[string]interface{})
	storepath, _ := os.Getwd()
	workdir := "/var/lib/mysql"
	storedir := filepath.Join(storepath, "mysql-storage")

	if strings.Compare(dbtype,"mongodb") == 0 {
		workdir = "/data/db"
		storedir = filepath.Join(storepath, "mongodb-storage")
	}
	if strings.Compare(dbtype,"mysql") == 0 {
		containerID, err = docker.CreateMysqlContainer(
			dbctx,
			cli,
			dockerImage,
			port,
			workdir,
			storedir,
			env)
	} else {
		containerID, err = docker.CreateMongoDBContainer(
			dbctx,
			cli,
			dockerImage,
			port,
			workdir,
			storedir,
			env)
	}
	err = docker.StartContainer(containerID)
	if err != nil {
		return "", types.NewResErr(500, "container not started", err)
	}

	return containerID, nil
}
