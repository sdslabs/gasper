package mongoDb

import (
	"github.com/sdslabs/SWS/lib/gin"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewServiceEngine()

func init() {
	Router.POST("/", validateRequest, createDB)
	Router.GET("/", fetchDBs)
	Router.GET("/logs", gin.FetchMysqlContainerLogs)
	Router.GET("/restart", gin.ReloadMysqlService)
	Router.GET("/db/:db", gin.FetchDBInfo)
	Router.DELETE("/", deleteDB)
}
