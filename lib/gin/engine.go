package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/middlewares"
)

// NewEngine returns a router setting up required configs
func NewEngine() *gin.Engine {
	if configs.GasperConfig.Debug {
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
