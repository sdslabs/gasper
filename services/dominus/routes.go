package dominus

import (
	g "github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/middlewares"
	"github.com/sdslabs/SWS/lib/utils"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewEngine()

func init() {
	Router.POST("/:service", checkFalconPlugIn(), createApp)
	Router.GET("/", checkFalconPlugIn(), fetchDocs)
	Router.GET("/:app", checkFalconPlugIn(), gin.FetchAppInfo)
	Router.GET("/:app/:action", checkFalconPlugIn(), execute)
	// Router.PUT("/", updateApp)
	// Router.DELETE("/", deleteApp)
}

func checkFalconPlugIn() g.HandlerFunc {
	if utils.FalconConfig["plugIn"].(bool) {
		return middlewares.FalconGuard()
	}
	return func(c *g.Context) {
		c.Next()
	}
}
