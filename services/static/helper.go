package static

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/api"
	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
)

type context struct {
	Index  string `json:"index" valid:"required~Field 'index' inside field 'context' was required but was not provided"`
	RcFile bool   `json:"rcFile"`
}

type staticRequestBody struct {
	Name    string                 `json:"name" valid:"required~Field 'name' is required but was not provided,alphanum~Field 'name' should only have alphanumeric characters,stringlength(3|40)~Field 'name' should have length between 3 to 40 characters"`
	URL     string                 `json:"url" valid:"required~Field 'url' is required but was not provided,url~Field 'url' is not a valid URL"`
	Context context                `json:"context"`
	Env     map[string]interface{} `json:"env"`
}

func validateRequest(c *gin.Context) {

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	// Restore the io.ReadCloser to its original state
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	var req staticRequestBody

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
