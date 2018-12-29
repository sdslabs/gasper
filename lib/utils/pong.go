package utils

import "github.com/gin-gonic/gin"

// Pong function is used for checking if a service is alive or not
func Pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
