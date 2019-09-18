package python

import (
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/sdslabs/SWS/lib/commons"
	"github.com/sdslabs/SWS/lib/configs"
	g "github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

// createApp function handles requests for making making new static app
func createApp(c *gin.Context) {
	var data map[string]interface{}
	c.BindJSON(&data)

	data["language"] = ServiceName
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
		utils.HostIP+configs.ServiceConfig[ServiceName].(map[string]interface{})["port"].(string),
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
		ServiceName,
		utils.HostIP+configs.ServiceConfig[ServiceName].(map[string]interface{})["port"].(string),
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

func rebuildApp(c *gin.Context) {
	appName := c.Param("app")
	filter := map[string]interface{}{
		"name":         appName,
		"language":     ServiceName,
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
