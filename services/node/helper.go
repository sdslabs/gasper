package node

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/api"
	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
)

type context struct {
	Index  string `json:"index" valid:"required~Field 'index' inside field 'context' was required but was not provided"`
	Port   string `json:"port" valid:"required~Field 'port' inside field 'context' was required but was not provided,port~Field 'port' inside field 'context' is not a valid port"`
	RcFile bool   `json:"rcFile"`
}

type nodeRequestBody struct {
	Name           string                 `json:"name" valid:"required~Field 'name' is required but was not provided,alphanum~Field 'name' should only have alphanumeric characters,stringlength(3|40)~Field 'name' should have length between 3 to 40 characters"`
	URL            string                 `json:"url" valid:"required~Field 'url' is required but was not provided,url~Field 'url' is not a valid URL"`
	Context        context                `json:"context"`
	NPM            bool                   `json:"npm"`
	Env            map[string]interface{} `json:"env"`
	GitAccessToken string                 `json:"git_access_token"`
}

func validateRequest(c *gin.Context) {

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	// Restore the io.ReadCloser to its original state
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	var req nodeRequestBody

	err := json.Unmarshal(bodyBytes, &req)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	if result, err := validator.ValidateStruct(req); !result {
		c.AbortWithStatusJSON(400, gin.H{
			"error": strings.Split(err.Error(), ";"),
		})
	} else {
		c.Next()
	}
}

// installPackages function installs the dependancies for the app
func installPackages(appEnv *types.ApplicationEnv) (string, types.ResponseError) {
	cmd := []string{"bash", "-c", `export NVM_DIR="$HOME/.nvm" && [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"; npm install &> /proc/1/fd/1`}
	execID, err := docker.ExecDetachedProcess(appEnv.Context, appEnv.Client, appEnv.ContainerID, cmd)
	if err != nil {
		return "", types.NewResErr(500, "Failed to perform npm install in the container", err)
	}
	return execID, nil
}

// startApp function starts the app using pm2
func startApp(index string, appEnv *types.ApplicationEnv) (string, types.ResponseError) {
	cmd := []string{"bash", "-c", `export NVM_DIR="$HOME/.nvm" && [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"; pm2 start ` + index + ` &> /proc/1/fd/1`}
	execID, err := docker.ExecDetachedProcess(appEnv.Context, appEnv.Client, appEnv.ContainerID, cmd)
	if err != nil {
		return "", types.NewResErr(500, "Failed to perform start app in the container", err)
	}
	return execID, nil
}

func pipeline(data map[string]interface{}) types.ResponseError {
	appConf := &types.ApplicationConfig{
		DockerImage:  utils.ServiceConfig["node"].(map[string]interface{})["image"].(string),
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

	var execID string
	// Perform npm install in the container
	if data["npm"] != nil {
		if data["npm"].(bool) {
			execID, resErr = installPackages(appEnv)
			if resErr != nil {
				go utils.FullCleanup(data["name"].(string))
				return resErr
			}
			data["execID"] = execID
		}
	}

	index := context["index"].(string)

	// Start app using pm2 in the container
	execID, resErr = startApp(index, appEnv)
	if resErr != nil {
		go utils.FullCleanup(data["name"].(string))
		return resErr
	}
	data["execID"] = execID

	return nil
}
