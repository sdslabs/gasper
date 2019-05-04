package php

import (
	"github.com/gin-gonic/gin"
	g "github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

// createApp function handles requests for making making new php app
func createApp(c *gin.Context) {
	var data map[string]interface{}
	c.BindJSON(&data)

	data["language"] = "php"

	resErr := pipeline(data)
	if resErr != nil {
		g.SendResponse(c, resErr, gin.H{})
		return
	}

	documentID, err := mongo.RegisterApp(data)

	if err != nil {
		go utils.FullCleanup(data["name"].(string))
		go utils.StateCleanup(data["name"].(string))
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	err = redis.RegisterApp(
		data["name"].(string),
		utils.HostIP+utils.ServiceConfig["php"].(map[string]interface{})["port"].(string),
	)

	if err != nil {
		go utils.FullCleanup(data["name"].(string))
		go utils.StateCleanup(data["name"].(string))
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	err = redis.IncrementServiceLoad(
		"php",
		utils.HostIP+utils.ServiceConfig["php"].(map[string]interface{})["port"].(string),
	)

	if err != nil {
		go utils.FullCleanup(data["name"].(string))
		go utils.StateCleanup(data["name"].(string))
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"id":      documentID,
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

func updateAppInfo(c *gin.Context) {
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

func rebuildApp(c *gin.Context) {
	appName := c.Param("app")
	filter := map[string]interface{}{
		"name":     appName,
		"language": "php",
	}
	data := mongo.FetchAppInfo(filter)[0]
	utils.FullCleanup(appName)

	resErr := pipeline(data)
	if resErr != nil {
		g.SendResponse(c, resErr, gin.H{})
		return
	}

	c.JSON(200, gin.H{
		"message": mongo.UpdateApp(filter, data),
	})
}
