package database

import (
	"fmt"
	"os"
	"path/filepath"

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
	fmt.Println("ck1")
	if err != nil {
		return "", types.NewResErr(500, "cannot setup client", err)
	}

<<<<<<< HEAD
	dockerImage := configs.ServiceConfig["mysql"].(map[string]interface{})["image"].(string)
	port := configs.ServiceConfig["mysql"].(map[string]interface{})["container_port"].(string)
	env := configs.ServiceConfig["mysql"].(map[string]interface{})["env"].(map[string]interface{})

=======
	dockerImage := utils.ServiceConfig[dbtype].(map[string]interface{})["image"].(string)
	port := utils.ServiceConfig[dbtype].(map[string]interface{})["container_port"].(string)
	env := utils.ServiceConfig[dbtype].(map[string]interface{})["env"].(map[string]interface{})
	fmt.Println("ck2")
>>>>>>> checked working
	storepath, _ := os.Getwd()
	workdir := "/var/lib/mysql"
	storedir := filepath.Join(storepath, "mysql-storage")

	if dbtype == "mongoDb" {
		workdir = "/var/lib/mongodb"
		storedir = filepath.Join(storepath, "mongodb-storage")
	}
	if dbtype == "mysql" {
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
<<<<<<< HEAD

	err = docker.StartContainer(containerID)
=======
	fmt.Println("ck3")
	err = docker.StartContainer(dbctx, cli, containerID)
>>>>>>> checked working
	if err != nil {
		return "", types.NewResErr(500, "container not started", err)
	}

	return containerID, nil
}
