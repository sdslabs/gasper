package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/factory"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/appmaker"
	"github.com/sdslabs/gasper/services/master/middlewares"
	"github.com/sdslabs/gasper/types"
)

type metricsRecord struct {
	UptimeRecord []bool    `json:"uptime_record"`
	CPURecord    []float64 `json:"cpu_record"`
	MemoryRecord []float64 `json:"memory_record"`
}

// FetchAppsByUser returns all applications owned by a user
func FetchAppsByUser(c *gin.Context) {
	fetchInstancesByUser(c, mongo.AppInstance)
}

// GetAllApplications gets all the applications from DB
func GetAllApplications(c *gin.Context) {
	fetchInstances(c, mongo.AppInstance)
}

// GetApplicationInfo gets info regarding a particular application
func GetApplicationInfo(c *gin.Context) {
	app := c.Param("app")
	filter := make(types.M)
	filter[mongo.NameKey] = app
	c.JSON(200, gin.H{
		"success": true,
		"data":    mongo.FetchAppInfo(filter),
	})
}

// BulkUpdateApps updates multiple application documents in mongoDB
func BulkUpdateApps(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)
	filter[mongo.InstanceTypeKey] = mongo.AppInstance

	var data types.M
	c.BindJSON(&data)

	err := validateUpdatePayload(data)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	_, err = mongo.UpdateInstances(filter, data)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	data["success"] = true
	c.JSON(200, data)
}

// UpdateAppByName updates the app getting name from url params
func UpdateAppByName(c *gin.Context) {
	app := c.Param("app")
	filter := types.M{
		mongo.NameKey:         app,
		mongo.InstanceTypeKey: mongo.AppInstance,
	}
	var data types.M
	c.BindJSON(&data)

	err := validateUpdatePayload(data)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	err = mongo.UpdateInstance(filter, data)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}

// CreateApp creates an application via gRPC
func CreateApp(c *gin.Context) {
	instanceURL, err := redis.GetLeastLoadedWorker()
	if err != nil {
		utils.SendServerErrorResponse(c, err)
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

	response, err := factory.CreateApplication(c.Param("language"), claims.GetEmail(), instanceURL, data)
	if err != nil {
		utils.LogError("Master-Controller-Application-1", err)
		if strings.Contains(err.Error(), "authentication required") {
			c.AbortWithStatusJSON(400, gin.H{
				"success": false,
				"error":   "Invalid git repository url or access token",
			})
		} else if strings.Contains(err.Error(), "couldn't find remote ref") {
			c.AbortWithStatusJSON(400, gin.H{
				"success": false,
				"error":   "Invalid git branch provided",
			})
		} else {
			utils.SendServerErrorResponse(c, err)
		}
		return
	}
	c.Data(200, "application/json", response)
}

// DeleteApp deletes an application via gRPC
func DeleteApp(c *gin.Context) {
	appName := c.Param("app")
	instanceURL, err := redis.FetchAppNode(appName)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Application %s is not deployed at the moment", appName),
		})
		return
	}

	response, err := factory.DeleteApplication(appName, instanceURL)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, response)
}

// FetchAppLogs returns the docker container logs of an application via gRPC
func FetchAppLogs(c *gin.Context) {
	appName := c.Param("app")
	instanceURL, err := redis.FetchAppNode(appName)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Application %s is not deployed at the moment", appName),
		})
		return
	}

	filter := utils.QueryToFilter(c.Request.URL.Query())
	if filter["tail"] == nil {
		filter["tail"] = "-1"
	}

	response, err := factory.FetchApplicationLogs(appName, filter["tail"].(string), instanceURL)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, response)
}

// RebuildApp rebuilds an application via gRPC
func RebuildApp(c *gin.Context) {
	appName := c.Param("app")
	instanceURL, err := redis.FetchAppNode(appName)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Application %s is not deployed at the moment", appName),
		})
		return
	}

	response, err := factory.RebuildApplication(appName, instanceURL)
	if err != nil {
		utils.LogError("Master-Controller-Application-2", err)
		if strings.Contains(err.Error(), "authentication required") {
			c.AbortWithStatusJSON(400, gin.H{
				"success": false,
				"error":   "Invalid git repository url or access token",
			})
		} else if strings.Contains(err.Error(), "couldn't find remote ref") {
			c.AbortWithStatusJSON(400, gin.H{
				"success": false,
				"error":   "Invalid git branch provided",
			})
		} else {
			utils.SendServerErrorResponse(c, err)
		}
		return
	}
	c.Data(200, "application/json", response)
}

// TransferApplicationOwnership transfers the ownership of an application to another user
func TransferApplicationOwnership(c *gin.Context) {
	transferOwnership(c, c.Param("app"), mongo.AppInstance, c.Param("user"))
}

// FetchMetrics retrieves the metrics of an application's container
func FetchMetrics(c *gin.Context) {
	appName := c.Param("app")
	filter := utils.QueryToFilter(c.Request.URL.Query())
	var timeSpan int64
	for unit, converter := range timeConversionMap {
		if val, ok := filter[unit].(string); ok {
			timeVal, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				continue
			}
			timeSpan += timeVal * converter
		}
	}

	metrics := mongo.FetchContainerMetrics(types.M{
		mongo.NameKey: appName,
		mongo.TimestampKey: types.M{
			"$gte": time.Now().Unix() - timeSpan,
		},
	}, -1)

	sparsity := timeConversionMap[filter["sparsity"].(string)]

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
		}

		CPURecord = append(CPURecord, metrics[i]["cpu_usage"].(float64))
		memoryRecord = append(memoryRecord, metrics[i]["memory_usage"].(float64))
	}

	metricsRecord := metricsRecord{uptimeRecord, CPURecord, memoryRecord}

	c.JSON(200, gin.H{
		"success": true,
		"data":    metricsRecord,
	})
}

func CreateRepository(c *gin.Context) {
	appName := c.Param("app")
	deployKey, err := appmaker.CreateRepository(appName)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   fmt.Sprintf(" %s github repository could not be created", appName),
		})
		return
	}

	c.JSON(200, gin.H{
		"success":   true,
		"deployKey": deployKey,
	})
}
