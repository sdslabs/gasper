package mysql

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type mysqlRequestBody struct {
	Name     string `json:"name" valid:"required~Field 'name' is required but was not provided,alphanum~Field 'name' should only have alphanumeric characters"`
	User     string `json:"user" valid:"required~Field 'user' is required but was not provided,alphanum~Field 'user' should only have alphanumeric characters"`
	Password string `json:"password" valid:"required~Field 'password' is required but was not provided,alphanum~Field 'password' should only have alphanumeric characters"`
}

func validateRequest(c *gin.Context) {

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	// Restore the io.ReadCloser to its original state
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	var req mysqlRequestBody

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
