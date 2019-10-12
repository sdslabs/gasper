package static

import (
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/api"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/types"
)

// Pipeline is the application creation pipeline
func Pipeline(data map[string]interface{}) types.ResponseError {
	appConf := &types.ApplicationConfig{
		ConfFunction: configs.CreateStaticContainerConfig,
		DockerImage:  configs.ImageConfig.Static,
	}

	appEnv, resErr := api.SetupApplication(appConf, data)
	if resErr != nil {
		return resErr
	}
	cmd := []string{"sh", "-c", `rm /etc/nginx/conf.d/default.conf && nginx -s reload`}
	_, err := docker.ExecDetachedProcess(appEnv.Context, appEnv.Client, appEnv.ContainerID, cmd)
	if err != nil {
		return types.NewResErr(500, "Failed to load application configuration", err)
	}
	return nil
}
