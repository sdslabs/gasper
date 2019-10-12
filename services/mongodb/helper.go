package mongodb

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/middlewares"
)

type mongodbRequestBody struct {
	Name     string `json:"name" valid:"required~Field 'name' is required but was not provided,alphanum~Field 'name' should only have alphanumeric characters,lowercase~Field 'name' should have only lowercase characters"`
	Password string `json:"password" valid:"required~Field 'password' is required but was not provided"`
}

// validateRequestBody validates the request body for the current microservice
func validateRequestBody(c *gin.Context) {
	middlewares.ValidateRequestBody(c, &mongodbRequestBody{})
}
