package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
)

// SendServerErrorResponse sends internal server error messages
// to the client depending on development mode or production mode
func SendServerErrorResponse(c *gin.Context, err error) {
	var errMessage string
	if configs.GasperConfig.Debug {
		errMessage = err.Error()
	} else {
		errMessage = "INTERNAL_SERVER_ERROR"
	}
	LogError(err)
	c.AbortWithStatusJSON(500, gin.H{
		"success": false,
		"error":   errMessage,
	})
}
