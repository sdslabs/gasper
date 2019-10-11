package mongodb

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/configs"
	"github.com/sdslabs/SWS/lib/commons"
	"github.com/sdslabs/SWS/lib/database"
	"github.com/sdslabs/SWS/lib/middlewares"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

func createDB(c *gin.Context) {
	userStr := middlewares.ExtractClaims(c)

	var data map[string]interface{}
	c.BindJSON(&data)

	delete(data, "rebuild")
	data["language"] = mongo.MongoDB
	data["instanceType"] = mongo.DBInstance
	data["hostIP"] = utils.HostIP
	data["containerPort"] = configs.ServiceConfig["mongodb"].(map[string]interface{})["container_port"].(string)
	data["owner"] = userStr.Email

	dbKey := fmt.Sprintf(`%s:%s`, data["user"].(string), data["name"].(string))

	err := database.CreateMongoDB(data["name"].(string), data["user"].(string), data["password"].(string))

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = mongo.UpsertInstance(
		map[string]interface{}{
			"name":         data["name"],
			"instanceType": data["instanceType"],
		}, data)

	if err != nil && err != mongo.ErrNoDocuments {
		go commons.DatabaseFullCleanup(dbKey, mongo.MongoDB)
		go commons.DatabaseStateCleanup(dbKey)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = redis.RegisterDB(
		dbKey,
		utils.HostIP+configs.ServiceConfig[ServiceName].(map[string]interface{})["port"].(string),
	)

	if err != nil {
		go commons.DatabaseFullCleanup(dbKey, mongo.MongoDB)
		go commons.DatabaseStateCleanup(dbKey)
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
		go commons.DatabaseFullCleanup(dbKey, mongo.MongoDB)
		go commons.DatabaseStateCleanup(dbKey)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	data["success"] = true
	c.JSON(200, data)
}

func fetchDBs(c *gin.Context) {
	userStr := middlewares.ExtractClaims(c)

	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = ServiceName
	filter["instanceType"] = mongo.DBInstance
	filter["owner"] = userStr.Email

	c.JSON(200, gin.H{
		"data": mongo.FetchDBs(filter),
	})
}

func deleteDB(c *gin.Context) {
	userStr := middlewares.ExtractClaims(c)

	user := c.Param("user")
	db := c.Param("db")
	dbKey := fmt.Sprintf(`%s:%s`, user, db)
	err := database.DeleteMongoDB(db, user)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = redis.RemoveDB(dbKey)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	filter := map[string]interface{}{
		"name":         db,
		"user":         user,
		"language":     ServiceName,
		"instanceType": mongo.MongoDB,
	}

	if !userStr.IsAdmin {
		filter["owner"] = userStr.Email
	}

	c.JSON(200, gin.H{
		"message": mongo.DeleteInstance(filter),
	})
}
