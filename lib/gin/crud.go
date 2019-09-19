package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/commons"
	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
)

// CreateApp creates an application for a service
func CreateApp(service string, pipeline func(data map[string]interface{}) types.ResponseError) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data map[string]interface{}
		c.BindJSON(&data)

		delete(data, "rebuild")
		data["language"] = service
		data["instanceType"] = mongo.AppInstance

		resErr := pipeline(data)
		if resErr != nil {
			SendResponse(c, resErr, gin.H{})
			return
		}

		err := mongo.UpsertInstance(
			map[string]interface{}{
				"name":         data["name"],
				"instanceType": data["instanceType"],
			}, data)

		if err != nil {
			go commons.FullCleanup(data["name"].(string), data["instanceType"].(string))
			go commons.StateCleanup(data["name"].(string), data["instanceType"].(string))
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}

		err = redis.RegisterApp(
			data["name"].(string),
			utils.HostIP+configs.ServiceConfig[service].(map[string]interface{})["port"].(string),
		)

		if err != nil {
			go commons.FullCleanup(data["name"].(string), data["instanceType"].(string))
			go commons.StateCleanup(data["name"].(string), data["instanceType"].(string))
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}

		err = redis.IncrementServiceLoad(
			service,
			utils.HostIP+configs.ServiceConfig[service].(map[string]interface{})["port"].(string),
		)

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
	}
}

// FetchAppInfo returns the information about a particular app
func FetchAppInfo(c *gin.Context) {
	app := c.Param("app")
	filter := make(map[string]interface{})
	filter["name"] = app
	filter["instanceType"] = mongo.AppInstance
	c.JSON(200, gin.H{
		"data": mongo.FetchAppInfo(filter),
	})
}

// FetchDBInfo returns the information about a particular db
func FetchDBInfo(c *gin.Context) {
	db := c.Param("db")
	filter := make(map[string]interface{})
	filter["name"] = db
	filter["instanceType"] = mongo.DBInstance
	c.JSON(200, gin.H{
		"data": mongo.FetchDBInfo(filter),
	})
}

// FetchDocs returns a handler function for fetching documents of all microservices of one kind
func FetchDocs(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		queries := c.Request.URL.Query()
		filter := utils.QueryToFilter(queries)
		if service != "dominus" {
			filter["language"] = service
		}
		c.JSON(200, gin.H{
			"data": mongo.FetchAppInfo(filter),
		})
	}
}

// DeleteApp returns a handler function for deleting an application bound to a microservice
func DeleteApp(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		queries := c.Request.URL.Query()
		filter := utils.QueryToFilter(queries)
		if service != "dominus" {
			filter["language"] = service
		}
		filter["instanceType"] = mongo.AppInstance
		c.JSON(200, gin.H{
			"message": mongo.DeleteInstance(filter),
		})
	}
}

// UpdateAppInfo returns a handler function for updating an application bound to a microservice
func UpdateAppInfo(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		queries := c.Request.URL.Query()
		filter := utils.QueryToFilter(queries)
		if service != "dominus" {
			filter["language"] = service
		}
		filter["instanceType"] = mongo.AppInstance
		var (
			data map[string]interface{}
		)
		c.BindJSON(&data)
		c.JSON(200, gin.H{
			"message": mongo.UpdateInstance(filter, data),
		})
	}
}
