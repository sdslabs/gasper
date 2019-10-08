package mysql

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/middlewares"
)

type mysqlRequestBody struct {
	Name string `json:"name" valid:"required~Field 'name' is required but was not provided,alphanum~Field 'name' should only have alphanumeric characters"`
	// User     string `json:"user" valid:"required~Field 'user' is required but was not provided,alphanum~Field 'user' should only have alphanumeric characters"`
	// Password string `json:"password" valid:"required~Field 'password' is required but was not provided,alphanum~Field 'password' should only have alphanumeric characters"`
}

// validateRequestBody validates the request body for the current microservice
func validateRequestBody(c *gin.Context) {
	middlewares.ValidateRequestBody(c, &mysqlRequestBody{})
}
