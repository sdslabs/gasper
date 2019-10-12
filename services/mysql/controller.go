package mysql

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/commons"
	"github.com/sdslabs/gasper/lib/database"
	"github.com/sdslabs/gasper/lib/middlewares"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
)

func createDB(c *gin.Context) {
	userStr := middlewares.ExtractClaims(c)

	var data map[string]interface{}
	c.BindJSON(&data)

	delete(data, "rebuild")
	data["language"] = mongo.Mysql
	data["instanceType"] = mongo.DBInstance
	data["hostIP"] = utils.HostIP
	data["containerPort"] = configs.ServiceConfig.Mysql.ContainerPort
	data["owner"] = userStr.Email

	data["user"] = data["name"].(string)

	db := data["name"].(string)

	if db == "root" {
		c.JSON(400, gin.H{
			"error": "Database name cannot be `root`",
		})
		return
	}
	err := database.CreateMysqlDB(data["name"].(string), data["user"].(string), data["password"].(string))
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
		go commons.DatabaseFullCleanup(db, mongo.Mysql)
		go commons.DatabaseStateCleanup(db)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = redis.RegisterDB(
		db,
		fmt.Sprintf("%s:%d", utils.HostIP, configs.ServiceConfig.Mysql.Port),
	)

	if err != nil {
		go commons.DatabaseFullCleanup(db, mongo.Mysql)
		go commons.DatabaseStateCleanup(db)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = redis.IncrementServiceLoad(
		ServiceName,
		fmt.Sprintf("%s:%d", utils.HostIP, configs.ServiceConfig.Mysql.Port),
	)

	if err != nil {
		go commons.DatabaseFullCleanup(db, mongo.Mysql)
		go commons.DatabaseStateCleanup(db)
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

	db := c.Param("db")

	err := database.DeleteMysqlDB(db)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = redis.RemoveDB(db)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	filter := map[string]interface{}{
		"name":         db,
		"language":     ServiceName,
		"instanceType": mongo.DBInstance,
	}

	if !userStr.IsAdmin {
		filter["owner"] = userStr.Email
	}

	c.JSON(200, gin.H{
		"message": mongo.DeleteInstance(filter),
	})
}
