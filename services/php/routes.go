package php

import (
	"github.com/sdslabs/SWS/lib/gin"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewServiceEngine()

// The parameters required for the creation of an application
var requiredParams = []string{"name", "url", "context", "composerPath", "composer"}

func init() {
	Router.POST("/", gin.ValidateParams(requiredParams), createApp)
	Router.GET("/", fetchDocs)
	Router.GET("/:app", gin.FetchAppInfo)
	Router.GET("/:app/logs", gin.FetchLogs)
	Router.GET("/:app/restart", gin.ReloadServer)
	Router.GET("/:app/rebuild", rebuildApp)
	Router.PUT("/", updateAppInfo)
	Router.DELETE("/", deleteApp)
}
