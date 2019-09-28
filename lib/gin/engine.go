package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/configs"
	"github.com/sdslabs/SWS/lib/middlewares"
)

// NewEngine returns a router setting up required configs
func NewEngine() *gin.Engine {
	if configs.SWSConfig["debug"].(bool) {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()
	return engine
}

// NewServiceEngine returns a router setting up required configs for micro-services
func NewServiceEngine() *gin.Engine {
	serviceEngine := NewEngine()
	serviceEngine.Use(middlewares.AuthorizeService())
	return serviceEngine
}
