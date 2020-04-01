package controllers

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/factory"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/kaze/middlewares"
	"github.com/sdslabs/gasper/types"
)

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
		utils.LogError(err)
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
		utils.LogError(err)
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
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	appName := c.Param("app")
	if appName == "" {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "App name not provided for metrics",
		})
		return
	}

	metricsInterval, err := strconv.ParseInt(c.Query("interval"), 10, 64)
	if err != nil {
		metricsInterval = 2 * int64(configs.ServiceConfig.Mizu.MetricsInterval)
	}

	metricsCount, err := strconv.ParseInt(c.Query("count"), 10, 64)
	if err != nil {
		metricsCount = 10
	}

	chanStream := make(chan []types.M, 10)
	go func() {
		defer close(chanStream)
		metrics := mongo.FetchContainerMetrics(types.M{
			mongo.NameKey: appName,
			mongo.TimestampKey: types.M{
				"$gte": time.Now().Unix() - int64(configs.ServiceConfig.Mizu.MetricsInterval*time.Second),
			},
		}, metricsCount)
		chanStream <- metrics
		if metricsInterval < int64(configs.ServiceConfig.Mizu.MetricsInterval) {
			metricsInterval = 2 * int64(configs.ServiceConfig.Mizu.MetricsInterval)
		}

		time.Sleep(time.Second * time.Duration(metricsInterval))
	}()

	c.Stream(func(w io.Writer) bool {
		if metrics, ok := <-chanStream; ok {
			c.SSEvent("metrics", metrics)
			return true
		}

		return false
	})
}
