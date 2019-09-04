package mongoDb

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/database"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

func createDB(c *gin.Context) {
	var data map[string]interface{}
	c.BindJSON(&data)

	data["language"] = "mongoDb"
	data["instanceType"] = mongo.DBInstance

	var dbKey = fmt.Sprintf(`%s:%s`, data["name"].(string), data["user"].(string))

	err := database.CreateMongoDB(data["dbname"].(string), data["dbuser"].(string), data["dbpass"].(string))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	}

	databaseID, err := mongo.RegisterInstance(data)

	if err != nil {
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
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	err = redis.IncrementServiceLoad(
		"mongoDb",
		utils.HostIP+utils.ServiceConfig["mongoDb"].(map[string]interface{})["port"].(string),
	)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"id":      databaseID,
	})
}

func fetchDBs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = "mongoDb"
	filter["instanceType"] = mongo.DBInstance

	c.JSON(200, gin.H{
		"data": mongo.FetchDBs(filter),
	})
}

func deleteDB(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	err := database.DeleteDB(filter["dbname"].(string), filter["dbuser"].(string))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	}

	filter["language"] = "mongoDb"
	filter["instanceType"] = mongo.DBInstance

	c.JSON(200, gin.H{
		"message": mongo.DeleteInstance(filter),
	})
}

func deleteMongoDB(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	err := database.DeleteMongoDB(filter["dbname"].(string), filter["dbuser"].(string), filter["dbpass"].(string))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	}

	filter["language"] = "mongoDb"
	filter["instanceType"] = mongo.DBInstance

	c.JSON(200, gin.H{
		"message": mongo.DeleteInstance(filter),
	})
}
