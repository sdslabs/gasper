package static

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/utils"
)

// Router is the main routes handler for the current microservice package
var Router = gin.Default()

func init() {
	Router.POST("/", createApp)
	Router.GET("/", fetchDocs)
	Router.GET("/ping", utils.Pong)
	Router.PUT("/", updateApp)
	Router.DELETE("/", deleteApp)
}
