package mongoDb

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/commons"
	"github.com/sdslabs/SWS/lib/database"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

func createDB(c *gin.Context) {
	var data map[string]interface{}
	c.BindJSON(&data)

	data["language"] = "mongoDb"
	data["instanceType"] = mongo.MongoDBInstance

	dbKey := fmt.Sprintf(`%s:%s:%s`, data["name"].(string), data["user"].(string), data["password"].(string))
	fmt.Println("check1")
	err := database.CreateMongoDB(data["name"].(string), data["user"].(string), data["password"].(string))
	fmt.Println("check2")
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	}

	databaseID, err := mongo.RegisterInstance(data)
	fmt.Println("check3")
	if err != nil {
		commons.FullCleanup(dbKey, data["instanceType"].(string))
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	err = redis.RegisterDB(
		dbKey,
		utils.HostIP+utils.ServiceConfig["mongoDb"].(map[string]interface{})["port"].(string),
	)

	if err != nil {
		commons.FullCleanup(dbKey, data["instanceType"].(string))
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	err = redis.IncrementServiceLoad(
		"mongoDb",
		utils.HostIP+utils.ServiceConfig["mongoDb"].(map[string]interface{})["port"].(string),
	)
	fmt.Println("check4")
	if err != nil {
		commons.FullCleanup(dbKey, data["instanceType"].(string))
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}
	fmt.Println("check 5")
	c.JSON(200, gin.H{
		"success": true,
		"id":      databaseID,
	})
}

func fetchDBs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = "mongoDb"
	filter["instanceType"] = mongo.MongoDBInstance

	c.JSON(200, gin.H{
		"data": mongo.FetchDBs(filter),
	})
}

func deleteDB(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	err := database.DeleteMongoDB(filter["name"].(string), filter["user"].(string), filter["password"].(string))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	}

	filter["language"] = "mongoDb"
	filter["instanceType"] = mongo.MongoDBInstance

	c.JSON(200, gin.H{
		"message": mongo.DeleteInstance(filter),
	})
}
