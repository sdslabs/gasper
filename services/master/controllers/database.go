package controllers

import (
	"errors"
    "fmt"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/factory"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/master/middlewares"
	"github.com/sdslabs/gasper/types"
)

// FetchDatabasesByUser returns all databases owned by a user
func FetchDatabasesByUser(c *gin.Context) {
	fetchInstancesByUser(c, mongo.DBInstance)
}

// GetAllDatabases gets all the Databases info from mongoDB
func GetAllDatabases(c *gin.Context) {
	fetchInstances(c, mongo.DBInstance)
}

// GetDatabaseInfo gets info regarding a particular database
func GetDatabaseInfo(c *gin.Context) {
	db := c.Param("db")
	filter := make(types.M)
	filter[mongo.NameKey] = db
	c.JSON(200, gin.H{
		"success": true,
		"data":    mongo.FetchDBInfo(filter),
	})
}

// CreateDatabase creates a database via gRPC
func CreateDatabase(c *gin.Context) {
	database := c.Param("database")
	instanceURL, err := redis.GetLeastLoadedInstance(database)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if instanceURL == redis.ErrEmptySet {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "No worker instances available at the moment",
		})
		return
	}

	data, err := c.GetRawData()
	if err != nil {
		utils.SendServerErrorResponse(c, errors.New("Failed to extract data from Request Body"))
		return
	}

	claims := middlewares.ExtractClaims(c)
	if claims == nil {
		utils.SendServerErrorResponse(c, errors.New("Failed to extract JWT claims"))
		return
	}

	response, err := factory.CreateDatabase(database, claims.GetEmail(), instanceURL, data)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.Data(200, "application/json", response)
}

// DeleteDatabase deletes a database via gRPC
func DeleteDatabase(c *gin.Context) {
	db := c.Param("db")
	instanceURL, err := redis.FetchDbNode(db)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "No such database exists",
		})
		return
	}

	response, err := factory.DeleteDatabase(db, instanceURL)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, response)
}

// FetchDatabaseLogs returns the docker container logs of a Database via gRPC
func FetchDatabaseLogs(c *gin.Context) {
	db := c.Param("db")
	instanceURL, err := redis.FetchDbNode(db)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Database %s is not deployed at the moment", db),
		})
		return
	}

	filter := utils.QueryToFilter(c.Request.URL.Query())
	if filter["tail"] == nil {
		filter["tail"] = "-1"
	}

	response, err := factory.FetchDatabaseServerLogs(db, filter["tail"].(string), instanceURL)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, response)
}

// TransferDatabaseOwnership transfers the ownership of a database to another user
func TransferDatabaseOwnership(c *gin.Context) {
	transferOwnership(c, c.Param("db"), mongo.DBInstance, c.Param("user"))
}
