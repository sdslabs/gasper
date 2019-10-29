package mizu

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/middlewares"
	"github.com/sdslabs/gasper/types"
)

var disallowedNames = []string{
	types.Dominus,
	types.Mizu,
	types.Hikari,
	types.Enrai,
	types.EnraiSSL,
	types.MySQL,
	types.MongoDB,
	types.SSH,
}

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
