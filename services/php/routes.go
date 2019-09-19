package php

import (
	"github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/middlewares"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewServiceEngine()

// ServiceName is the name of the current microservice
var ServiceName = "php"

func init() {
	Router.POST("/", validateRequestBody, middlewares.IsUniqueApp(), gin.CreateApp(ServiceName, pipeline))
	Router.GET("/", gin.FetchDocs(ServiceName))
	Router.GET("/:app", gin.FetchAppInfo)
	Router.GET("/:app/logs", gin.FetchLogs)
	Router.GET("/:app/restart", gin.ReloadServer)
	Router.GET("/:app/rebuild", rebuildApp)
	Router.PUT("/", gin.UpdateAppInfo(ServiceName))
	Router.DELETE("/", gin.DeleteApp(ServiceName))
}
