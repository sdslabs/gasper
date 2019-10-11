package mizu

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/configs"
	"github.com/sdslabs/SWS/lib/cloudflare"
	"github.com/sdslabs/SWS/lib/commons"
	g "github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/middlewares"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// createApp creates an application for a given language
func createApp(c *gin.Context) {
	userStr := middlewares.ExtractClaims(c)

	language := c.Param("language")
	var data map[string]interface{}
	c.BindJSON(&data)

	delete(data, "rebuild")
	data["language"] = language
	data["instanceType"] = mongo.AppInstance
	data["owner"] = userStr.Email

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
		fmt.Sprintf("%s:%d", utils.HostIP, configs.ServiceConfig.Mizu.Port),
		fmt.Sprintf("%s:%d", utils.HostIP, data["httpPort"].(int)),
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
		fmt.Sprintf("%s:%d", utils.HostIP, configs.ServiceConfig.Mizu.Port),
	)

	if err != nil {
		go commons.AppFullCleanup(data["name"].(string))
		go commons.AppStateCleanup(data["name"].(string))
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	if configs.CloudflareConfig.PlugIn {
		resp, err := cloudflare.CreateRecord(data["name"].(string), mongo.AppInstance)
		if err != nil {
			go commons.AppFullCleanup(data["name"].(string))
			go commons.AppStateCleanup(data["name"].(string))
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		data["cloudflareID"] = resp.Result.ID
		data["domainURL"] = fmt.Sprintf("%s.%s.%s", data["name"].(string), mongo.AppInstance, configs.GasperConfig.Domain)
	}

	data["success"] = true
	c.JSON(200, data)
}

func rebuildApp(c *gin.Context) {
	userStr := middlewares.ExtractClaims(c)

	appName := c.Param("app")
	filter := map[string]interface{}{
		"name":         appName,
		"instanceType": mongo.AppInstance,
		"owner":        userStr.Email,
	}
	dataList := mongo.FetchAppInfo(filter)
	if len(dataList) == 0 {
		c.JSON(400, gin.H{
			"error": "No such application exists",
		})
		return
	}
	data := dataList[0]
	data["context"] = map[string]interface{}(data["context"].(primitive.M))
	data["resources"] = map[string]interface{}(data["resources"].(primitive.M))

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
	userStr := middlewares.ExtractClaims(c)

	app := c.Param("app")
	filter := map[string]interface{}{
		"name":         app,
		"instanceType": mongo.AppInstance,
	}

	if !userStr.IsAdmin {
		filter["owner"] = userStr.Email
	}
	update := map[string]interface{}{
		"deleted": true,
	}
	node, _ := redis.FetchAppNode(app)
	go redis.DecrementServiceLoad(ServiceName, node)
	go redis.RemoveApp(app)
	go commons.AppFullCleanup(app)
	if configs.CloudflareConfig.PlugIn {
		go cloudflare.DeleteRecord(app, mongo.AppInstance)
	}
	c.JSON(200, gin.H{
		"message": mongo.UpdateInstance(filter, update),
	})
}
