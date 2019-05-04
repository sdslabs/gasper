package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/middlewares"
)

// NewEngine returns a router setting up required configs
func NewEngine() *gin.Engine {
	if configs.SWSConfig["debug"].(bool) {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	return gin.Default()
}

// NewServiceEngine returns a router setting up required configs for micro-services
func NewServiceEngine() *gin.Engine {
	engine := NewEngine()
	engine.Use(middlewares.AuthorizeService())
	return engine
}
