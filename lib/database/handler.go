package database

import (
	"github.com/docker/docker/client"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
	"golang.org/x/net/context"
)

// SetupDBInstance sets up the mysql instance for deployment
func SetupDBInstance() (string, types.ResponseError) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", types.NewResErr(500, "cannot setup client", err)
	}

	dockerImage := utils.ServiceConfig["mysql"].(map[string]interface{})["image"].(string)
	port := utils.ServiceConfig["mysql"].(map[string]interface{})["container_port"].(string)
	env := utils.ServiceConfig["mysql"].(map[string]interface{})["env"].(map[string]interface{})
	workdir := "/var/lib/mysql"
	storedir := "mysql-storage"

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
