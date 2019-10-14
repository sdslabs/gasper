package mongodb

import (
	"github.com/sdslabs/gasper/lib/gin"
	"github.com/sdslabs/gasper/lib/middlewares"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewServiceEngine()

// ServiceName is the name of the current microservice
var ServiceName = "mongodb"

func init() {
	Router.POST("/mongodb", validateRequestBody, middlewares.IsUniqueDB(), createDB)
	Router.GET("", fetchDBs)
	Router.GET("/logs", gin.FetchMongoDBContainerLogs)
	Router.GET("/restart", gin.ReloadMongoDBService)
	Router.GET("/db/:db", gin.FetchDBInfo)
	Router.DELETE("/:db", deleteDB)
}
