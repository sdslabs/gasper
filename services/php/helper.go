package php

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/api"
	"github.com/sdslabs/SWS/lib/commons"
	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/middlewares"
	"github.com/sdslabs/SWS/lib/types"
)

type context struct {
	Index  string `json:"index" valid:"required~Field 'index' inside field 'context' was required but was not provided"`
	RcFile bool   `json:"rcFile"`
}

type phpRequestBody struct {
	Name           string                     `json:"name" valid:"required~Field 'name' is required but was not provided,alphanum~Field 'name' should only have alphanumeric characters,stringlength(3|40)~Field 'name' should have length between 3 to 40 characters,lowercase~Field 'name' should have only lowercase characters"`
	Password       string                     `json:"password" valid:"required~Field 'password' is required but was not provided,alphanum~Field 'password' should only have alphanumeric characters"`
	URL            string                     `json:"url" valid:"required~Field 'url' is required but was not provided,url~Field 'url' is not a valid URL"`
	Context        context                    `json:"context"`
	Resources      types.ApplicationResources `json:"resources"`
	Composer       bool                       `json:"composer"`
	ComposerPath   string                     `json:"composerPath"`
	Env            map[string]interface{}     `json:"env"`
	GitAccessToken string                     `json:"git_access_token"`
}

// validateRequestBody validates the request body for the current microservice
func validateRequestBody(c *gin.Context) {
	middlewares.ValidateRequestBody(c, &phpRequestBody{})
}

// installPackages installs dependancies for the specific microservice
func installPackages(path string, appEnv *types.ApplicationEnv) (string, types.ResponseError) {
	cmd := []string{"bash", "-c", `composer install -d ` + path}
	execID, err := docker.ExecDetachedProcess(appEnv.Context, appEnv.Client, appEnv.ContainerID, cmd)
	if err != nil {
		return "", types.NewResErr(500, "Failed to perform composer install in the container", err)
	}
	return execID, nil
}

func pipeline(data map[string]interface{}) types.ResponseError {
	appConf := &types.ApplicationConfig{
		DockerImage:  configs.ServiceConfig["php"].(map[string]interface{})["image"].(string),
		ConfFunction: configs.CreatePHPContainerConfig,
	}

	appEnv, resErr := api.SetupApplication(appConf, data)
	if resErr != nil {
		return resErr
	}

	context := data["context"].(map[string]interface{})

	if context["rcFile"].(bool) {
		return nil
	}

	// Perform composer install in the container
	if data["composer"] != nil {
		if data["composer"].(bool) {
			var composerPath string
			if data["composerPath"] != nil {
				composerPath = data["composerPath"].(string)
			} else {
				composerPath = "."
			}
			execID, resErr := installPackages(composerPath, appEnv)
			if resErr != nil {
				go commons.FullCleanup(data["name"].(string), data["instanceType"].(string))
				return resErr
			}
			data["execID"] = execID
		}
	}

	return nil
}
