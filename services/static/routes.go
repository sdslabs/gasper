package static

import (
	"github.com/gin-gonic/gin"
)

// Router is the main routes handler for the current microservice package
var Router = gin.Default()

func init() {
	appGroup := Router.Group("/static")
	{
		appGroup.POST("/", createApp)
		appGroup.GET("/", fetchDocs)
		appGroup.PUT("/", updateApp)
		appGroup.DELETE("/", deleteApp)
	}
}
