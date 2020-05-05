package database

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/sdslabs/gasper/lib/utils"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/types"
)

// CreateRedisDBContainer  creates a postgre database
func CreateRedisDBContainer(db types.Database) error {
	storepath, _ := os.Getwd()
	var err error
	dockerImage := configs.ImageConfig.Redis
	port, err := utils.GetFreePort()

	contaierport := strconv.Itoa(port)
	env := configs.ServiceConfig.Kaen.RedisKaen.Env
	workdir := "/data/"
	storedir := filepath.Join(storepath, "kaen-redis-storage", db.GetName())
	_, err = docker.CreateRedisContainer(
		dockerImage,
		contaierport,
		workdir,
		storedir,
		db.GetName(),
		env)
	if err != nil {
		return types.NewResErr(500, "container not created", err)
	}
	return nil
}

// DeleteRedisDBContainer delete container
func DeleteRedisDBContainer(containerID string) error {
	err := docker.DeleteContainer(containerID)
	return err

}
