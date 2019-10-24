package mizu

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/cloudflare"
	"github.com/sdslabs/gasper/lib/commons"
	g "github.com/sdslabs/gasper/lib/gin"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// createApp creates an application for a given language
func createApp(c *gin.Context) {
	language := c.Param("language")
	var data types.M
	c.BindJSON(&data)

	delete(data, "rebuild")
	data["language"] = language
	data["instanceType"] = mongo.AppInstance

	resErr := pipeline(data)
	if resErr != nil {
		g.SendResponse(c, resErr, gin.H{})
		return
	}

	if configs.CloudflareConfig.PlugIn {
		resp, err := cloudflare.CreateRecord(data["name"].(string), mongo.AppInstance)
		if err != nil {
			go commons.AppFullCleanup(data["name"].(string))
			go commons.AppStateCleanup(data["name"].(string))
			utils.SendServerErrorResponse(c, err)
			return
		}
		data["cloudflareID"] = resp.Result.ID
		data["domainURL"] = fmt.Sprintf("%s.%s.%s", data["name"].(string), mongo.AppInstance, configs.GasperConfig.Domain)
	}

	err := mongo.UpsertInstance(
		types.M{
			"name":         data["name"],
			"instanceType": data["instanceType"],
		}, data)

	if err != nil && err != mongo.ErrNoDocuments {
		go commons.AppFullCleanup(data["name"].(string))
		go commons.AppStateCleanup(data["name"].(string))
		utils.SendServerErrorResponse(c, err)
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
		utils.SendServerErrorResponse(c, err)
		return
	}

	err = redis.IncrementServiceLoad(
		ServiceName,
		fmt.Sprintf("%s:%d", utils.HostIP, configs.ServiceConfig.Mizu.Port),
	)

	if err != nil {
		go commons.AppFullCleanup(data["name"].(string))
		go commons.AppStateCleanup(data["name"].(string))
		utils.SendServerErrorResponse(c, err)
		return
	}

	data["success"] = true
	c.JSON(200, data)
}

func rebuildApp(c *gin.Context) {
	appName := c.Param("app")
	filter := types.M{
		"name": appName,
	}

	dataList := mongo.FetchAppInfo(filter)
	if len(dataList) == 0 {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "No such application exists",
		})
		return
	}
	data := dataList[0]
	data["context"] = types.M(data["context"].(primitive.M))
	data["resources"] = types.M(data["resources"].(primitive.M))

	commons.AppFullCleanup(appName)

	if componentMap[data["language"].(string)] == nil {
		utils.SendServerErrorResponse(c, fmt.Errorf("Non-supported language `%s` specified for `%s`", data["language"].(string), appName))
		return
	}
	resErr := pipeline(data)
	if resErr != nil {
		fmt.Println("error mila 1")
		g.SendResponse(c, resErr, gin.H{})
		return
	}
	fmt.Println("error nhi mila 1")
	err := mongo.UpdateInstance(filter, data)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	fmt.Println("error nhi mila 2")
	c.JSON(200, gin.H{
		"success": true,
	})
}

// deleteApp deletes an application in a worker node
func deleteApp(c *gin.Context) {
	app := c.Param("app")
	filter := types.M{
		"name":         app,
		"instanceType": mongo.AppInstance,
	}

	node, _ := redis.FetchAppNode(app)
	go redis.DecrementServiceLoad(ServiceName, node)
	go redis.RemoveApp(app)
	go commons.AppFullCleanup(app)
	if configs.CloudflareConfig.PlugIn {
		go cloudflare.DeleteRecord(app, mongo.AppInstance)
	}

	_, err := mongo.DeleteInstance(filter)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}
