package mizu

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/middlewares"
	"github.com/sdslabs/gasper/types"
)

func validateRequestBody(c *gin.Context) {
	language := c.Param("language")
	if pipeline[language] == nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Language `%s` is not supported", language),
		})
		return
	}
	middlewares.ValidateRequestBody(c, &types.ApplicationConfig{})
}
