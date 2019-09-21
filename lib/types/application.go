package types

import (
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// AppBinding defines the struct for storing both the server and node urls
type AppBinding struct {
	Node   string `json:"node"`
	Server string `json:"server"`
}

// ApplicationConfig defines the config for various service apps
type ApplicationConfig struct {
	DockerImage  string
	ConfFunction func(string, map[string]interface{}) string
}

// ApplicationEnv defines the environment of the running app
type ApplicationEnv struct {
	Context     context.Context
	Client      *client.Client
	ContainerID string
}

// ApplicationResources defines the resources requested by an app
type ApplicationResources struct {
	// Memory limits in GB
	Memory float64 `json:"memory" valid:"float~Field 'memory' inside field 'resources' should be of type float"`

	// CPU quota in units of CPUs
	CPU float64 `json:"cpu" valid:"float~Field 'cpu' inside field 'resources' should be of type float"`
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
