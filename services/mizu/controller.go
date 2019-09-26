package mizu

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/sdslabs/SWS/configs"
	"github.com/sdslabs/SWS/lib/cloudflare"
	"github.com/sdslabs/SWS/lib/commons"
	g "github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

// createApp creates an application for a given language
func createApp(c *gin.Context) {
	language := c.Param("language")
	var data map[string]interface{}
	c.BindJSON(&data)

	delete(data, "rebuild")
	data["language"] = language
	data["instanceType"] = mongo.AppInstance

	resErr := componentMap[language].pipeline(data)
	if resErr != nil {
		g.SendResponse(c, resErr, gin.H{})
		return
	}

	err := mongo.UpsertInstance(
		map[string]interface{}{
			"name":         data["name"],
			"instanceType": data["instanceType"],
		}, data)

	if err != nil && err != mongo.ErrNoDocuments {
		go commons.AppFullCleanup(data["name"].(string))
		go commons.AppStateCleanup(data["name"].(string))
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = redis.RegisterApp(
		data["name"].(string),
		utils.HostIP+configs.ServiceConfig[ServiceName].(map[string]interface{})["port"].(string),
		utils.HostIP+":"+strconv.Itoa(data["httpPort"].(int)),
	)

	if err != nil {
		go commons.AppFullCleanup(data["name"].(string))
		go commons.AppStateCleanup(data["name"].(string))
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = redis.IncrementServiceLoad(
		ServiceName,
		utils.HostIP+configs.ServiceConfig[ServiceName].(map[string]interface{})["port"].(string),
	)

	if err != nil {
		go commons.AppFullCleanup(data["name"].(string))
		go commons.AppStateCleanup(data["name"].(string))
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	if configs.CloudflareConfig["plugIn"].(bool) {
		resp, err := cloudflare.CreateRecord(data["name"].(string), mongo.AppInstance, utils.HostIP)
		if err != nil {
			go commons.AppFullCleanup(data["name"].(string))
			go commons.AppStateCleanup(data["name"].(string))
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		data["cloudflareID"] = resp.Result.ID
		data["domainURL"] = fmt.Sprintf("%s.%s.%s", data["name"].(string), mongo.AppInstance, configs.SWSConfig["domain"].(string))
	}

	data["success"] = true
	c.JSON(200, data)
}

func rebuildApp(c *gin.Context) {
	appName := c.Param("app")
	filter := map[string]interface{}{
		"name":         appName,
		"instanceType": mongo.AppInstance,
	}
	dataList := mongo.FetchAppInfo(filter)
	if len(dataList) == 0 {
		c.JSON(400, gin.H{
			"error": "No such application exists",
		})
		return
	}
	data := dataList[0]
	data["context"] = map[string]interface{}(data["context"].(primitive.D).Map())

	commons.AppFullCleanup(appName)

	resErr := componentMap[data["language"].(string)].pipeline(data)
	if resErr != nil {
		g.SendResponse(c, resErr, gin.H{})
		return
	}

	c.JSON(200, gin.H{
		"message": mongo.UpdateInstance(filter, data),
	})
}

// deleteApp deletes an application in a worker node
func deleteApp(c *gin.Context) {
	app := c.Param("app")
	filter := map[string]interface{}{
		"name":         app,
		"instanceType": mongo.AppInstance,
	}
	update := map[string]interface{}{
		"deleted": true,
	}
	node, _ := redis.FetchAppNode(app)
	go redis.DecrementServiceLoad(ServiceName, node)
	go redis.RemoveApp(app)
	go commons.AppFullCleanup(app)
	if configs.CloudflareConfig["plugIn"].(bool) {
		go cloudflare.DeleteRecord(app, mongo.AppInstance)
	}
	c.JSON(200, gin.H{
		"message": mongo.UpdateInstance(filter, update),
	})
}
