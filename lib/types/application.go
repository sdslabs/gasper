package types

import (
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// ApplicationConfig defines the config for various service apps
type ApplicationConfig struct {
	DockerImage  string
	ConfFunction func(string, string) string
}

// ApplicationEnv defines the environment of the running app
type ApplicationEnv struct {
	Context     context.Context
	Client      *client.Client
	ContainerID string
}

// NewAppEnv returns a new ApplicationEnv
func NewAppEnv() (*ApplicationEnv, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &ApplicationEnv{
		Context: ctx,
		Client:  cli,
	}, nil
}

// StaticAppConfig defined the request structure for creating new static app
type StaticAppConfig struct {
	Name      string `json:"name" form:"name" binding:"required"`
	UserID    int    `json:"user_id" form:"user_id" binding:"required"`
	GithubURL string `json:"github_url" form:"github_url" binding:"required"`
}
