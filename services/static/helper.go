package static

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/api"
	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
)

type context struct {
	Index string `json:"index" valid:"required"`
}

type staticRequestBody struct {
	Name    string                 `json:"name" valid:"required,alphanum,stringlength(3|40)"`
	URL     string                 `json:"url" valid:"required,url"`
	Context context                `json:"context" valid:"required"`
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

	json.Unmarshal(bodyBytes, &req)

	if result, err := validator.ValidateStruct(req); !result {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err,
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
