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

	err := database.CreateDB(data["dbname"].(string), data["dbuser"].(string), data["dbpass"].(string))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	}

	databaseID, err := mongo.RegisterDB(data)

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

	c.JSON(200, gin.H{
		"data": mongo.FetchDBs(filter),
	})
}

func updateDB(c *gin.Context) {
	// function to update a database by importing the schema
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

	c.JSON(200, gin.H{
		"message": mongo.DeleteDB(filter),
	})
}
