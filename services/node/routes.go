package node

import (
	"github.com/sdslabs/SWS/lib/gin"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewServiceEngine()

func init() {
	Router.POST("/", createApp)
	Router.GET("/", fetchDocs)
	Router.GET("/:app", gin.FetchAppInfo)
	Router.GET("/:app/logs", gin.FetchLogs)
	Router.GET("/:app/restart", gin.ReloadServer)
	Router.PUT("/", updateAppInfo)
	Router.DELETE("/", deleteApp)
}
