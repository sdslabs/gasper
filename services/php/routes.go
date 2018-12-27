package php

import (
	"github.com/gin-gonic/gin"
)

// Router is the main routes handler for the current microservice package
var Router = gin.Default()

func init() {
	Router.POST("/", createApp)
	Router.GET("/", fetchDocs)
	Router.PUT("/", updateApp)
	Router.DELETE("/", deleteApp)
}
