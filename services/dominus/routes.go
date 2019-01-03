package dominus

import (
	"github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/utils"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewEngine()

func init() {
	Router.POST("/:service", createApp)
	Router.GET("/", fetchDocs)
	Router.GET("/ping", utils.Pong)
	// Router.PUT("/", updateApp)
	// Router.DELETE("/", deleteApp)
}
