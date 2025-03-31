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
	"github.com/sdslabs/gasper/services/master/middlewares"
	"github.com/sdslabs/gasper/types"
)

type metricsRecord struct {
	UptimeRecord   bool    `json:"uptime_record"`
	TimeStamp      int64   `json:"time_stamp"`
	CPUUsage       float64 `json:"cpu_usage"`
	MemoryUsage    float64 `json:"memory_usage"`
	HostIP         string  `json:"host_ip"`
	OnlineCPUs     float64 `json:"online_cpus"`
	MaxMemoryUsage float64 `json:"max_memory_usage"`
	MemoryLimit    float64 `json:"memory_limit"`
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
	appName := c.Param("app")
	filter := types.M{
		mongo.NameKey:         appName,
		mongo.InstanceTypeKey: mongo.AppInstance,
	}
	var data types.UpdatePayload
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	instanceURL, err := redis.FetchAppNode(appName)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Application %s is not deployed at the moment", appName),
		})
		return
	}
	app, err := mongo.FetchSingleApp(appName)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}
	err = UpdateData(app, &data)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	err = mongo.UpdateInstance(filter, app)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	response, err := factory.UpdateApplication(appName, instanceURL)
	if err != nil {
		utils.LogError("Master-Controller-Application-3", err)
		if strings.Contains(err.Error(), "authentication required") {
			c.AbortWithStatusJSON(400, gin.H{
				"success": false,
				"error":   "Invalid git repository url or access token",
			})
		} else if strings.Contains(err.Error(), "invalid reference") {
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

// DeleteAppUsingAppName deletes an application via gRPC using appName as a parameter
func DeleteAppUsingAppname(c *gin.Context, appName string) {
	instanceURL, err := redis.FetchAppNode(appName)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Application %s is not deployed at the moment", appName),
		})
		return
	}

	_, err = factory.DeleteApplication(appName, instanceURL)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
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
	val := filter["tail"]
	if val == nil || val == "" {
		val = "-1"
	} else if _, err := strconv.Atoi(val.(string)); err != nil {
		utils.SendServerErrorResponse(c, errors.New("invalid tail value "+val.(string)))
		return
	}

	response, err := factory.FetchApplicationLogs(appName, val.(string), instanceURL)
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

	metrics := mongo.FetchContainerMetrics(types.M{
		mongo.NameKey: appName,
		mongo.TimestampKey: types.M{
			"$gte": time.Now().Unix() - timeSpan,
		},
	}, -1)

	// handle empty metrics
	if len(metrics) == 0 {
		c.JSON(200, types.M{
			"success": true,
			"data":    []types.M{},
		})
		return
	}

	if val, ok := filter["sparsityvalue"].(string); ok {
		sparsityVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			sparsity = sparsityVal * timeConversionMap[filter["sparsityunit"].(string)]
		}
	}

	baseTimestamp := metrics[0]["timestamp"].(int64)
	var downtimeIntensity int = 0
	var currTimestamp int64

	records := []metricsRecord{}
	for i := range metrics {
		currTimestamp = metrics[i]["timestamp"].(int64)
		metricsRecord := metricsRecord{}
		if !metrics[i]["alive"].(bool) {
			downtimeIntensity++
		}
		if (baseTimestamp - currTimestamp) >= sparsity {
			//

			baseTimestamp = currTimestamp
			if downtimeIntensity > 0 {
				//
				metricsRecord.UptimeRecord = false
			} else {
				//
				metricsRecord.UptimeRecord = false
			}
			downtimeIntensity = 0
			metricsRecord.CPUUsage = metrics[i]["cpu_usage"].(float64)
			metricsRecord.MemoryUsage = metrics[i]["memory_usage"].(float64)
			metricsRecord.TimeStamp = metrics[i]["timestamp"].(int64)
			metricsRecord.HostIP = metrics[i]["host_ip"].(string)
			metricsRecord.OnlineCPUs = metrics[i]["online_cpus"].(float64)
			metricsRecord.MaxMemoryUsage = metrics[i]["max_memory_usage"].(float64)
			metricsRecord.MemoryLimit = metrics[i]["memory_limit"].(float64)
		}
		records = append(records, metricsRecord)
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    records,
	})
}
