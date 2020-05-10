package docker

import (
	"os"

	"github.com/docker/docker/client"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/utils"
	"golang.org/x/net/context"
)

// NewClient returns a new docker client
func NewClient() *client.Client {
	// No need to test docker connection if AppMaker and DbMaker are not deployed
	if !configs.ServiceConfig.AppMaker.Deploy && !configs.ServiceConfig.DbMaker.Deploy {
		return nil
	}
	cli, err := client.NewEnvClient()
	if err != nil {
		utils.Log("Failed creating Docker Client", utils.ErrorTAG)
		utils.LogError(err)
		os.Exit(1)
	}
	_, err = cli.Ping(context.Background())
	if err != nil {
		utils.Log("Connection with Docker Daemon was not established", utils.ErrorTAG)
		utils.LogError(err)
		os.Exit(1)
	}
	utils.LogInfo("Docker Daemon Connection Established")
	return cli
}

var cli = NewClient()
