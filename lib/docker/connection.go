package docker

import (
	"github.com/docker/docker/client"
	"github.com/sdslabs/gasper/lib/utils"
	"golang.org/x/net/context"
)

// NewClient returns a new docker client
func NewClient() *client.Client {
	cli, err := client.NewEnvClient()
	if err != nil {
		utils.Log("Failed creating Docker Client", utils.ErrorTAG)
		utils.LogError(err)
		panic(err)
	}
	_, err = cli.Ping(context.Background())
	if err != nil {
		utils.Log("Connection with Docker Daemon was not established", utils.ErrorTAG)
		utils.LogError(err)
		panic(err)
	}
	utils.LogInfo("Docker Daemon Connection Established")
	return cli
}

var cli = NewClient()