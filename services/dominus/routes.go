package dominus

import (
	"github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/middlewares"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewEngine()

func init() {
	Router.POST("/:service", middlewares.FalconGuard(), createApp)
	Router.GET("/", middlewares.FalconGuard(), fetchDocs)
	Router.GET("/:app", middlewares.FalconGuard(), gin.FetchAppInfo)
	Router.GET("/:app/:action", middlewares.FalconGuard(), execute)
	// Router.PUT("/", updateApp)
	// Router.DELETE("/", deleteApp)
}
