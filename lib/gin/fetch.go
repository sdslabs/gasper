package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/mongo"
)

// FetchAppInfo returns the information about a particular app
func FetchAppInfo(c *gin.Context) {
	app := c.Param("app")

	filter := make(map[string]interface{})

	filter["name"] = app

	c.JSON(200, gin.H{
		"data": mongo.FetchAppInfo(filter),
	})
}

// FetchDBInfo returns the information about a particular db
func FetchDBInfo(c *gin.Context) {
	db := c.Param("db")

	filter := make(map[string]interface{})

	filter["name"] = db

	c.JSON(200, gin.H{
		"data": mongo.FetchDBInfo(filter),
	})
}
