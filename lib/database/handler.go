package database

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
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

	dockerImage := utils.ServiceConfig[dbtype].(map[string]interface{})["image"].(string)
	port := utils.ServiceConfig[dbtype].(map[string]interface{})["container_port"].(string)
	env := utils.ServiceConfig[dbtype].(map[string]interface{})["env"].(map[string]interface{})

	storepath, _ := os.Getwd()
	workdir := "/var/lib/mysql"
	storedir := filepath.Join(storepath, "mysql-storage")

	if dbtype == "mongodb" {
		workdir = "/var/lib/mongodb"
		storedir = filepath.Join(storepath, "mongodb-storage")
	}
	containerID, err = docker.CreateMysqlContainer(
		dbctx,
		cli,
		dockerImage,
		port,
		workdir,
		storedir,
		env)

	err = docker.StartContainer(dbctx, cli, containerID)
	if err != nil {
		return "", types.NewResErr(500, "container not started", err)
	}

	return containerID, nil
}
