package static

import (
	"github.com/sdslabs/SWS/lib/gin"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewEngine()

func init() {
	Router.POST("/", createApp)
	Router.GET("/", fetchDocs)
	Router.GET("/:app", gin.FetchAppInfo)
	Router.GET("/:app/logs", gin.FetchLogs)
	Router.PUT("/", updateApp)
	Router.DELETE("/", deleteApp)
}
