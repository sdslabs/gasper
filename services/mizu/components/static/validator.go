package static

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/middlewares"
	"github.com/sdslabs/gasper/types"
)

type context struct {
	Index  string `json:"index" valid:"required~Field 'index' inside field 'context' was required but was not provided"`
	RcFile bool   `json:"rcFile"`
}

type staticRequestBody struct {
	Name           string                     `json:"name" valid:"required~Field 'name' is required but was not provided,alphanum~Field 'name' should only have alphanumeric characters,stringlength(3|40)~Field 'name' should have length between 3 to 40 characters,lowercase~Field 'name' should have only lowercase characters"`
	Password       string                     `json:"password" valid:"required~Field 'password' is required but was not provided"`
	URL            string                     `json:"url" valid:"required~Field 'url' is required but was not provided,url~Field 'url' is not a valid URL"`
	Context        context                    `json:"context"`
	Resources      types.ApplicationResources `json:"resources"`
	Env            map[string]interface{}     `json:"env"`
	GitAccessToken string                     `json:"git_access_token"`
}

// Validator validates the request body for static applications
func Validator(c *gin.Context) {
	middlewares.ValidateRequestBody(c, &staticRequestBody{})
}
