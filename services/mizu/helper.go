package mizu

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func validateRequestBody(c *gin.Context) {
	language := c.Param("language")
	if componentMap[language] == nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": fmt.Sprintf("Language `%s` is not supported", language),
		})
		return
	}
	componentMap[language].validator(c)
}
