package mizu

import (
	"net/http"

	"github.com/sdslabs/gasper/lib/gin"
	"github.com/sdslabs/gasper/types"
)

// ServiceName is the name of the current microservice
const ServiceName = types.Mizu

// NewService returns a new instance of the current microservice
func NewService() http.Handler {
	// router is the main routes handler for the current microservice package
	router := gin.NewServiceEngine()

	router.POST("/:language", createApp)
	router.GET("", gin.FetchDocs)
	router.GET("/:app", gin.FetchAppInfo)
	router.GET("/:app/logs", gin.FetchLogs)
	router.PATCH("/:app/rebuild", rebuildApp)
	router.PUT("/:app", gin.UpdateAppInfo)
	router.DELETE("/:app", deleteApp)

	return router
}
