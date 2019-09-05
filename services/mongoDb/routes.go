package mongoDb

import (
	"github.com/sdslabs/SWS/lib/gin"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewServiceEngine()

func init() {
	Router.POST("/", validateRequest, createDB)
	Router.GET("/", fetchDBs)
	Router.GET("/logs", gin.FetchMongoDBContainerLogs)
	Router.GET("/restart", gin.ReloadMongoDBService)
	Router.GET("/db/:db", gin.FetchDBInfo)
	Router.DELETE("/", deleteDB)
}
