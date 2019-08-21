package mysql

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/database"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/utils"
)

func createDB(c *gin.Context) {
	var data map[string]interface{}
	c.BindJSON(&data)

	data["language"] = "mysql"
	data["instanceType"] = mongo.DBInstance

	err := database.CreateDB(data["dbname"].(string), data["dbuser"].(string), data["dbpass"].(string))
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

	c.JSON(200, gin.H{
		"success": true,
		"id":      databaseID,
	})
}

func fetchDBs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = "mysql"
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

	filter["language"] = "mysql"
	filter["instanceType"] = mongo.DBInstance

	c.JSON(200, gin.H{
		"message": mongo.DeleteInstance(filter),
	})
}
