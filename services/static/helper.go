package static

import (
	"github.com/sdslabs/SWS/lib/api"
	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
)

func pipeline(data map[string]interface{}) types.ResponseError {
	appConf := &types.ApplicationConfig{
		ConfFunction: configs.CreateStaticContainerConfig,
		DockerImage:  utils.ServiceConfig["static"].(map[string]interface{})["image"].(string),
	}

	_, resErr := api.SetupApplication(appConf, data)
	if resErr != nil {
		return resErr
	}

	return nil
}
