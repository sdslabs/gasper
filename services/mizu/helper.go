package mizu

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/api"
	"github.com/sdslabs/gasper/types"
)

const (
	pythonVersionTag = "python_version"
	python3Tag       = "3"
	python2Tag       = "2"
)

func validateRequestBody(c *gin.Context) {
	language := c.Param("language")
	if componentMap[language] == nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Language `%s` is not supported", language),
		})
		return
	}
	componentMap[language].validator(c)
}

func pipeline(data types.M) types.ResponseError {
	var appConf types.ApplicationConfig
	switch data["language"].(string) {

	case types.Nodejs:
		appConf.DockerImage = configs.ImageConfig.Nodejs
		appConf.ConfFunction = configs.CreateNodeContainerConfig

	case types.Php:
		appConf.DockerImage = configs.ImageConfig.Php
		appConf.ConfFunction = configs.CreatePHPContainerConfig

	case types.Python:
		var image string
		if data[pythonVersionTag].(string) == python3Tag {
			image = configs.ImageConfig.Python3
		} else if data[pythonVersionTag].(string) == python2Tag {
			image = configs.ImageConfig.Python2
		}
		appConf.DockerImage = image
		appConf.ConfFunction = configs.CreatePythonContainerConfig

	case types.Static:
		appConf.DockerImage = configs.ImageConfig.Static
		appConf.ConfFunction = configs.CreateStaticContainerConfig
	}

	_, resErr := api.SetupApplication(&appConf, data)
	if resErr != nil {
		return resErr
	}

	return nil
}
