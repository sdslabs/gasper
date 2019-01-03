package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/utils"
)

// NewEngine returns a router setting up required configs
func NewEngine() *gin.Engine {
	if utils.SWSConfig["debug"].(bool) {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	return gin.Default()
}
