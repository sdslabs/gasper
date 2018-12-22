package php

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/mongo"
)

// createApp function handles requests for making making new php app
func createApp(c *gin.Context) {
	var (
		data map[string]interface{}
	)
	c.BindJSON(&data)
	data["language"] = "php"

	c.JSON(200, gin.H{
		"success": true,
		"id":      mongo.RegisterApp(data),
	})
}

func fetchDocs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := queryToFilter(queries)

	filter["language"] = "php"

	c.JSON(200, gin.H{
		"data": mongo.FetchAppInfo(filter),
	})
}

func deleteApp(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := queryToFilter(queries)

	filter["language"] = "php"

	c.JSON(200, gin.H{
		"message": mongo.DeleteApp(filter),
	})
}

func updateApp(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := queryToFilter(queries)

	filter["language"] = "php"

	var (
		data map[string]interface{}
	)
	c.BindJSON(&data)

	c.JSON(200, gin.H{
		"message": mongo.UpdateApp(filter, data),
	})
}

func queryToFilter(queries map[string][]string) map[string]interface{} {
	filter := make(map[string]interface{})

	for key, value := range queries {
		filter[key] = value[0]
	}

	return filter
}
