package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/configs"
)

// AuthorizeService creates a gin middleware to authorize dominus requests
func AuthorizeService() gin.HandlerFunc {
	secret := configs.SWSConfig["secret"].(string)
	return func(c *gin.Context) {
		dominusSecret := c.GetHeader("dominus-secret")
		if dominusSecret == "" {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "Missing 'dominus-secret' header",
			})
			return
		}
		if dominusSecret != secret {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Invalid 'dominus-secret'",
			})
			return
		}
		c.Next()
	}
}
