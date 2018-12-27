package php

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/utils"
)

// createApp function handles requests for making making new php app
func createApp(c *gin.Context) {
	var (
		data map[string]interface{}
	)

	c.BindJSON(&data)
	data["language"] = "php"

	var composerPath string = data["composerPath"].(string)

	// Perform compeser install in the container
	if data["composer"] == "true" {
		execId, err := installPackages(composerPath)

		// TODO: use execId and err later, for now just printing it out
		fmt.Println(execId)
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"success": true,
		"id":      mongo.RegisterApp(data),
	})
}

func fetchDocs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = "php"

	c.JSON(200, gin.H{
		"data": mongo.FetchAppInfo(filter),
	})
}

func deleteApp(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = "php"

	c.JSON(200, gin.H{
		"message": mongo.DeleteApp(filter),
	})
}

func updateApp(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = "php"

	var (
		data map[string]interface{}
	)
	c.BindJSON(&data)

	c.JSON(200, gin.H{
		"message": mongo.UpdateApp(filter, data),
	})
}
