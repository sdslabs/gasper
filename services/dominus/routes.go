package dominus

import (
	"github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/middlewares"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewEngine()

// ServiceName is the name of the current microservice
var ServiceName = "dominus"

func init() {
	Router.Use(middlewares.FalconGuard())
	Router.POST("/:service", createApp)
	Router.GET("/", gin.FetchDocs(ServiceName))
	Router.PUT("/", gin.UpdateAppInfo(ServiceName))
	Router.DELETE("/", gin.DeleteApp(ServiceName))
	app := Router.Group("/app")
	{
		app.GET("/:app", gin.FetchAppInfo)
		app.GET("/:app/:action", trimURLPath, execute)
	}
	db := Router.Group("/db")
	{
		db.GET("/:db", gin.FetchDBInfo)

	}
}
