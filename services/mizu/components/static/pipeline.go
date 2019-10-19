package static

import (
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/api"
	"github.com/sdslabs/gasper/types"
)

// Pipeline is the application creation pipeline
func Pipeline(data map[string]interface{}) types.ResponseError {
	appConf := &types.ApplicationConfig{
		ConfFunction: configs.CreateStaticContainerConfig,
		DockerImage:  configs.ImageConfig.Static,
	}

	_, resErr := api.SetupApplication(appConf, data)
	if resErr != nil {
		return resErr
	}
	return nil
}
