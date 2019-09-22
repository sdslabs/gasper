package mysql

import (
	"github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/middlewares"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewServiceEngine()

// ServiceName is the name of the current microservice
var ServiceName = "mysql"

func init() {
	Router.POST("/", validateRequestBody, middlewares.IsUniqueDB(), createDB)
	Router.GET("/", fetchDBs)
	Router.GET("/logs", gin.FetchMysqlContainerLogs)
	Router.GET("/restart", gin.ReloadMysqlService)
	Router.GET("/db/:db", gin.FetchDBInfo)
	Router.DELETE("/:user/:db", deleteDB)
}
