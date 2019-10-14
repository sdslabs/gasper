package mizu

import (
	"github.com/sdslabs/gasper/lib/gin"
	"github.com/sdslabs/gasper/lib/middlewares"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewServiceEngine()

// ServiceName is the name of the current microservice
var ServiceName = "mizu"

func init() {
	Router.POST("/:language", validateRequestBody, middlewares.IsUniqueApp(), createApp)
	Router.GET("", gin.FetchDocs)
	Router.GET("/:app", gin.FetchAppInfo)
	Router.GET("/:app/logs", gin.FetchLogs)
	Router.GET("/:app/restart", gin.ReloadServer)
	Router.GET("/:app/rebuild", rebuildApp)
	Router.PUT("/:app", gin.UpdateAppInfo)
	Router.DELETE("/:app", deleteApp)
}
