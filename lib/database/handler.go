package database

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
	"golang.org/x/net/context"
)

// SetupDBInstance sets up the mysql instance for deployment
func SetupDBInstance() (string, types.ResponseError) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", types.NewResErr(500, "cannot setup client", err)
	}

	dockerImage := configs.ServiceConfig["mysql"].(map[string]interface{})["image"].(string)
	port := configs.ServiceConfig["mysql"].(map[string]interface{})["container_port"].(string)
	env := configs.ServiceConfig["mysql"].(map[string]interface{})["env"].(map[string]interface{})

	storepath, _ := os.Getwd()
	workdir := "/var/lib/mysql"
	storedir := filepath.Join(storepath, "mysql-storage")

	containerID, err := docker.CreateMysqlContainer(
		ctx,
		cli,
		dockerImage,
		port,
		workdir,
		storedir,
		env)

	err = docker.StartContainer(ctx, cli, containerID)
	if err != nil {
		return "", types.NewResErr(500, "container not started", err)
	}

	return containerID, nil
}
