package dominus

import (
	"github.com/sdslabs/SWS/lib/gin"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewEngine()

func init() {
	Router.POST("/:service", createApp)
	Router.GET("/", fetchDocs)
	Router.GET("/:app", gin.FetchAppInfo)
	Router.GET("/:app/:action", execute)
	// Router.PUT("/", updateApp)
	// Router.DELETE("/", deleteApp)
}
