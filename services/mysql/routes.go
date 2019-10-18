package mysql

import (
	"net/http"

	"github.com/sdslabs/gasper/lib/gin"
	"github.com/sdslabs/gasper/lib/middlewares"
)

// ServiceName is the name of the current microservice
const ServiceName = "mysql"

// NewService returns a new instance of the current microservice
func NewService() http.Handler {
	// router is the main routes handler for the current microservice package
	router := gin.NewServiceEngine()

	router.POST("/mysql", validateRequestBody, middlewares.IsUniqueDB(), createDB)
	router.GET("", fetchDBs)
	router.GET("/logs", gin.FetchMysqlContainerLogs)
	router.GET("/restart", gin.ReloadMysqlService)
	router.GET("/db/:db", gin.FetchDBInfo)
	router.DELETE("/:db", deleteDB)

	return router
}
