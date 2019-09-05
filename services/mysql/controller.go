package mysql

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

	data["language"] = "mysql"
	data["instanceType"] = mongo.MysqlDBInstance

	var dbKey = fmt.Sprintf(`%s:%s`, data["name"].(string), data["user"].(string))
	fmt.Println("check1")
	err := database.CreateMysqlDB(data["name"].(string), data["user"].(string), data["password"].(string))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	}
	fmt.Println("check2")
	databaseID, err := mongo.RegisterInstance(data)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	err = redis.RegisterDB(
		dbKey,
		utils.HostIP+utils.ServiceConfig["mysql"].(map[string]interface{})["port"].(string),
	)
	fmt.Println("check3")
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	err = redis.IncrementServiceLoad(
		"mysql",
		utils.HostIP+utils.ServiceConfig["mysql"].(map[string]interface{})["port"].(string),
	)
	fmt.Println("check4")
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
	fmt.Println("check5")
}

func fetchDBs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = "mysql"
	filter["instanceType"] = mongo.MysqlDBInstance

	c.JSON(200, gin.H{
		"data": mongo.FetchDBs(filter),
	})
}

func deleteDB(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	err := database.DeleteMysqlDB(filter["name"].(string), filter["user"].(string))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	}

	filter["language"] = "mysql"
	filter["instanceType"] = mongo.MysqlDBInstance

	c.JSON(200, gin.H{
		"message": mongo.DeleteInstance(filter),
	})
}
