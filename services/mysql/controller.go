package mysql

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/commons"
	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/database"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

func createDB(c *gin.Context) {
	var data map[string]interface{}
	c.BindJSON(&data)

<<<<<<< HEAD
	delete(data, "rebuild")
	data["language"] = ServiceName
	data["instanceType"] = mongo.DBInstance
	data["hostIP"] = utils.HostIP
	data["containerPort"] = configs.ServiceConfig["mysql"].(map[string]interface{})["container_port"].(string)

	dbKey := fmt.Sprintf(`%s:%s`, data["user"].(string), data["name"].(string))

=======
	data["language"] = "mysql"
	data["instanceType"] = mongo.MysqlDBInstance

	var dbKey = fmt.Sprintf(`%s:%s`, data["name"].(string), data["user"].(string))
	fmt.Println("check1")
>>>>>>> checked working
	err := database.CreateMysqlDB(data["name"].(string), data["user"].(string), data["password"].(string))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}
<<<<<<< HEAD

	err = mongo.UpsertInstance(
		map[string]interface{}{
			"name":         data["name"],
			"instanceType": data["instanceType"],
		}, data)
=======
	fmt.Println("check2")
	databaseID, err := mongo.RegisterInstance(data)
>>>>>>> checked working

	if err != nil {
		go commons.FullCleanup(data["name"].(string), data["instanceType"].(string))
		go commons.StateCleanup(data["name"].(string), data["instanceType"].(string))
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	err = redis.RegisterDB(
		dbKey,
		utils.HostIP+configs.ServiceConfig[ServiceName].(map[string]interface{})["port"].(string),
	)
	fmt.Println("check3")
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
	fmt.Println("check4")
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
	})
	fmt.Println("check5")
}

func fetchDBs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

<<<<<<< HEAD
	filter["language"] = ServiceName
	filter["instanceType"] = mongo.DBInstance
=======
	filter["language"] = "mysql"
	filter["instanceType"] = mongo.MysqlDBInstance
>>>>>>> checked working

	c.JSON(200, gin.H{
		"data": mongo.FetchDBs(filter),
	})
}

func deleteDB(c *gin.Context) {
	user := c.Param("user")
	db := c.Param("db")
	dbKey := fmt.Sprintf(`%s:%s`, user, db)

	err := database.DeleteMysqlDB(db, user)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

<<<<<<< HEAD
	err = redis.RemoveDB(dbKey)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	filter := map[string]interface{}{
		"name":         db,
		"user":         user,
		"language":     ServiceName,
		"instanceType": mongo.DBInstance,
	}
=======
	filter["language"] = "mysql"
	filter["instanceType"] = mongo.MysqlDBInstance
>>>>>>> checked working

	c.JSON(200, gin.H{
		"message": mongo.DeleteInstance(filter),
	})
}
