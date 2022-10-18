package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"time"

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

// FetchMetrics retrieves the metrics of a Database's container
func FetchDatabaseMetrics(c *gin.Context) {
	db := c.Param("db")
	filter := utils.QueryToFilter(c.Request.URL.Query())
	var timeSpan int64
	var sparsity int64
	for unit, converter := range timeConversionMap {
		if val, ok := filter[unit].(string); ok {
			timeVal, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				continue
			}
			timeSpan += timeVal * converter
		}
	}

	language, _ := mongo.FetchDatabaseLanguage(db)

	metrics := mongo.FetchContainerMetrics(types.M{
		//temproary fix, currently showing logs of container rather than actual database
		mongo.NameKey: language,
		mongo.TimestampKey: types.M{
			"$gte": time.Now().Unix() - timeSpan,
		},
	}, -1)


	if val, ok := filter["sparsityvalue"].(string); ok {
		sparsityVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			sparsity = sparsityVal * timeConversionMap[filter["sparsityunit"].(string)]
		}
	}

	uptimeRecord := []bool{}
	CPURecord := []float64{}
	memoryRecord := []float64{}

	baseTimestamp := metrics[0]["timestamp"].(int64)
	var downtimeIntensity int = 0
	var currTimestamp int64

	for i := range metrics {
		currTimestamp = metrics[i]["timestamp"].(int64)
		if !metrics[i]["alive"].(bool) {
			downtimeIntensity++
		}
		if (baseTimestamp - currTimestamp) >= sparsity {
			baseTimestamp = currTimestamp
			if downtimeIntensity > 0 {
				uptimeRecord = append(uptimeRecord, false)
			} else {
				uptimeRecord = append(uptimeRecord, true)
			}
			downtimeIntensity = 0
			CPURecord = append(CPURecord, metrics[i]["cpu_usage"].(float64))
			memoryRecord = append(memoryRecord, metrics[i]["memory_usage"].(float64))
		}

	}

	metricsRecord := metricsRecord{uptimeRecord, CPURecord, memoryRecord}

	c.JSON(200, gin.H{
		"success": true,
		"data":    metricsRecord,
	})
}
