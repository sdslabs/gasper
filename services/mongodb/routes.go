package mongodb

import (
	"net/http"

	"github.com/sdslabs/gasper/lib/gin"
	"github.com/sdslabs/gasper/lib/middlewares"
)

// ServiceName is the name of the current microservice
const ServiceName = "mongodb"

// NewService returns a new instance of the current microservice
func NewService() http.Handler {
	// router is the main routes handler for the current microservice package
	router := gin.NewServiceEngine()

	router.POST("/mongodb", validateRequestBody, middlewares.IsUniqueDB(), createDB)
	router.GET("", fetchDBs)
	router.GET("/logs", gin.FetchMongoDBContainerLogs)
	router.GET("/restart", gin.ReloadMongoDBService)
	router.GET("/db/:db", gin.FetchDBInfo)
	router.DELETE("/:db", deleteDB)

	return router
}
