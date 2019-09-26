package python

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/middlewares"
	"github.com/sdslabs/SWS/lib/types"
)

type context struct {
	Index  string   `json:"index" valid:"required~Field 'index' inside field 'context' was required but was not provided"`
	Port   string   `json:"port" valid:"required~Field 'port' inside field 'context' was required but was not provided,port~Field 'port' inside field 'context' is not a valid port"`
	Args   []string `json:"args"`
	RcFile bool     `json:"rcFile"`
}

type pythonRequestBody struct {
	Name           string                     `json:"name" valid:"required~Field 'name' is required but was not provided,alphanum~Field 'name' should only have alphanumeric characters,stringlength(3|40)~Field 'name' should have length between 3 to 40 characters,lowercase~Field 'name' should have only lowercase characters"`
	Password       string                     `json:"password" valid:"required~Field 'password' is required but was not provided,alphanum~Field 'password' should only have alphanumeric characters"`
	URL            string                     `json:"url" valid:"required~Field 'url' is required but was not provided,url~Field 'url' is not a valid URL"`
	Context        context                    `json:"context"`
	Resources      types.ApplicationResources `json:"resources"`
	PythonVersion  string                     `json:"python_version" valid:"required~Field 'python_version' is required but was not provided"`
	Requirements   string                     `json:"requirements" valid:"required~Field 'requirements' is required but was not provided"`
	Django         bool                       `json:"django"`
	Env            map[string]interface{}     `json:"env"`
	GitAccessToken string                     `json:"git_access_token"`
}

// Validator validates the request body for python applications
func Validator(c *gin.Context) {
	middlewares.ValidateRequestBody(c, &pythonRequestBody{})
}
