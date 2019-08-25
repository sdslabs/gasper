package python

import (
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/sdslabs/SWS/lib/commons"
	g "github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

// createApp function handles requests for making making new static app
func createApp(c *gin.Context) {
	var data map[string]interface{}
	c.BindJSON(&data)

	data["language"] = "python"
	data["instanceType"] = mongo.AppInstance

	resErr := pipeline(data)
	if resErr != nil {
		g.SendResponse(c, resErr, gin.H{})
		return
	}

	documentID, err := mongo.RegisterInstance(data)

	if err != nil {
		go commons.FullCleanup(data["name"].(string), data["instanceType"].(string))
		go commons.StateCleanup(data["name"].(string), data["instanceType"].(string))
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	err = redis.RegisterApp(
		data["name"].(string),
		utils.HostIP+utils.ServiceConfig["python"].(map[string]interface{})["port"].(string),
	)

	if err != nil {
		go commons.FullCleanup(data["name"].(string), data["instanceType"].(string))
		go commons.StateCleanup(data["name"].(string), data["instanceType"].(string))
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	err = redis.IncrementServiceLoad(
		"python",
		utils.HostIP+utils.ServiceConfig["python"].(map[string]interface{})["port"].(string),
	)

	if err != nil {
		go commons.FullCleanup(data["name"].(string), data["instanceType"].(string))
		go commons.StateCleanup(data["name"].(string), data["instanceType"].(string))
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

	filter["language"] = "python"
	filter["instanceType"] = mongo.AppInstance

	c.JSON(200, gin.H{
		"data": mongo.FetchAppInfo(filter),
	})
}

func deleteApp(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = "python"
	filter["instanceType"] = mongo.AppInstance

	c.JSON(200, gin.H{
		"message": mongo.DeleteInstance(filter),
	})
}

func updateAppInfo(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = "python"
	filter["instanceType"] = mongo.AppInstance

	var (
		data map[string]interface{}
	)
	c.BindJSON(&data)

	c.JSON(200, gin.H{
		"message": mongo.UpdateInstance(filter, data),
	})
}

func rebuildApp(c *gin.Context) {
	appName := c.Param("app")
	filter := map[string]interface{}{
		"name":         appName,
		"language":     "python",
		"instanceType": mongo.AppInstance,
	}
	data := mongo.FetchAppInfo(filter)[0]
	data["context"] = map[string]interface{}(data["context"].(primitive.D).Map())

	commons.FullCleanup(appName, mongo.AppInstance)

	resErr := pipeline(data)
	if resErr != nil {
		g.SendResponse(c, resErr, gin.H{})
		return
	}

	c.JSON(200, gin.H{
		"message": mongo.UpdateInstance(filter, data),
	})
}
