package static

import (
	"github.com/sdslabs/SWS/configs"
	"github.com/sdslabs/SWS/lib/api"
	"github.com/sdslabs/SWS/lib/types"
)

// Pipeline is the application creation pipeline
func Pipeline(data map[string]interface{}) types.ResponseError {
	appConf := &types.ApplicationConfig{
		ConfFunction: configs.CreateStaticContainerConfig,
		DockerImage:  configs.ServiceConfig["static"].(map[string]interface{})["image"].(string),
	}

	_, resErr := api.SetupApplication(appConf, data)
	if resErr != nil {
		return resErr
	}

	return nil
}
