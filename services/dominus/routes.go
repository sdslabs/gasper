package dominus

import (
	"github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/middlewares"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewEngine()

func init() {
	Router.POST("/:service", middlewares.FalconGuard(), middlewares.CORSMiddleware(), createApp)
	Router.GET("/", middlewares.FalconGuard(), middlewares.CORSMiddleware(), fetchDocs)
	Router.GET("/:app", middlewares.FalconGuard(), middlewares.CORSMiddleware(), gin.FetchAppInfo)
	Router.GET("/:app/:action", middlewares.FalconGuard(), middlewares.CORSMiddleware(), execute)
	// Router.PUT("/", updateApp)
	// Router.DELETE("/", deleteApp)
}
