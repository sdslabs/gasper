package nodejs

import (
	"fmt"

	"github.com/sdslabs/SWS/configs"
	"github.com/sdslabs/SWS/lib/api"
	"github.com/sdslabs/SWS/lib/commons"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
)

// startApp function starts the app using pm2
func bootstrap(index string, appEnv *types.ApplicationEnv) (string, types.ResponseError) {
	cmd := []string{"sh", "-c", fmt.Sprintf(`npm install &> /proc/1/fd/1; node %s &> /proc/1/fd/1`, index)}
	execID, err := docker.ExecDetachedProcess(appEnv.Context, appEnv.Client, appEnv.ContainerID, cmd)
	if err != nil {
		return "", types.NewResErr(500, "Failed to launch application in the container", err)
	}
	return execID, nil
}

// Pipeline is the application creation pipeline
func Pipeline(data map[string]interface{}) types.ResponseError {
	appConf := &types.ApplicationConfig{
		DockerImage:  configs.ImageConfig.Nodejs,
		ConfFunction: configs.CreateNodeContainerConfig,
	}

	appEnv, resErr := api.SetupApplication(appConf, data)
	if resErr != nil {
		return resErr
	}

	context := data["context"].(map[string]interface{})

	if context["rcFile"].(bool) {
		return nil
	}

	index := context["index"].(string)

	// Start app using pm2 in the container
	execID, resErr := bootstrap(index, appEnv)
	if resErr != nil {
		go commons.AppFullCleanup(data["name"].(string))
		return resErr
	}
	data["execID"] = execID

	return nil
}
