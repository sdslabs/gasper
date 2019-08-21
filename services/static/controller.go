package static

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	g "github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

// createApp function handles requests for making making new static app
func createApp(c *gin.Context) {
	var data map[string]interface{}
	c.BindJSON(&data)

	data["language"] = "static"
	data["instanceType"] = "app"

	resErr := pipeline(data)
	if resErr != nil {
		g.SendResponse(c, resErr, gin.H{})
		return
	}

	documentID, err := mongo.RegisterInstance(data)

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
		utils.HostIP+utils.ServiceConfig["static"].(map[string]interface{})["port"].(string),
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
		"static",
		utils.HostIP+utils.ServiceConfig["static"].(map[string]interface{})["port"].(string),
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

	filter["language"] = "static"
	filter["instanceType"] = mongo.AppInstance

	c.JSON(200, gin.H{
		"data": mongo.FetchAppInfo(filter),
	})
}

func deleteApp(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = "static"
	filter["instanceType"] = mongo.AppInstance

	c.JSON(200, gin.H{
		"message": mongo.DeleteInstance(filter),
	})
}

func updateAppInfo(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = "static"
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
		"language":     "static",
		"instanceType": mongo.AppInstance,
	}
	data := mongo.FetchAppInfo(filter)[0]
	data["context"] = map[string]interface{}(data["context"].(primitive.D).Map())

	utils.FullCleanup(appName)

	resErr := pipeline(data)
	if resErr != nil {
		g.SendResponse(c, resErr, gin.H{})
		return
	}

	c.JSON(200, gin.H{
		"message": mongo.UpdateInstance(filter, data),
	})
}
