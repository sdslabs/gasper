package mongodb

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/commons"
	"github.com/sdslabs/gasper/lib/database"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

func createDB(c *gin.Context) {
	var data types.M
	c.BindJSON(&data)

	delete(data, "rebuild")
	data["language"] = mongo.MongoDB
	data["instanceType"] = mongo.DBInstance
	data["hostIP"] = utils.HostIP
	data["containerPort"] = configs.ServiceConfig.Mongodb.ContainerPort

	data["user"] = data["name"].(string)

	db := data["name"].(string)

	if db == "admin" {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "Database name cannot be `admin`",
		})
		return
	}
	err := database.CreateMongoDB(data["name"].(string), data["user"].(string), data["password"].(string))

	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	err = mongo.UpsertInstance(
		types.M{
			"name":         data["name"],
			"instanceType": data["instanceType"],
		}, data)

	if err != nil && err != mongo.ErrNoDocuments {
		go commons.DatabaseFullCleanup(db, mongo.MongoDB)
		go commons.DatabaseStateCleanup(db)
		utils.SendServerErrorResponse(c, err)
		return
	}

	err = redis.RegisterDB(
		db,
		fmt.Sprintf("%s:%d", utils.HostIP, configs.ServiceConfig.Mongodb.Port),
	)

	if err != nil {
		go commons.DatabaseFullCleanup(db, mongo.MongoDB)
		go commons.DatabaseStateCleanup(db)
		utils.SendServerErrorResponse(c, err)
		return
	}

	err = redis.IncrementServiceLoad(
		ServiceName,
		fmt.Sprintf("%s:%d", utils.HostIP, configs.ServiceConfig.Mongodb.Port),
	)

	if err != nil {
		go commons.DatabaseFullCleanup(db, mongo.MongoDB)
		go commons.DatabaseStateCleanup(db)
		utils.SendServerErrorResponse(c, err)
		return
	}
	data["success"] = true
	c.JSON(200, data)
}

func fetchDBs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = ServiceName
	filter["instanceType"] = mongo.DBInstance

	c.JSON(200, gin.H{
		"success": true,
		"data":    mongo.FetchDBs(filter),
	})
}

func deleteDB(c *gin.Context) {
	db := c.Param("db")
	err := database.DeleteMongoDB(db)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	err = redis.RemoveDB(db)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	filter := types.M{
		"name":         db,
		"language":     ServiceName,
		"instanceType": mongo.DBInstance,
	}

	_, err = mongo.DeleteInstance(filter)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}
