package gin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/commons"
	"github.com/sdslabs/SWS/configs"
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
			utils.HostIP+":"+strconv.Itoa(data["httpPort"].(int)),
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

// UpdateAppByName updates the app getting name from url params
func UpdateAppByName(c *gin.Context) {
	app := c.Param("app")
	filter := map[string]interface{}{
		"name":         app,
		"instanceType": mongo.AppInstance,
	}
	var data map[string]interface{}
	c.BindJSON(&data)
	c.JSON(200, gin.H{
		"message": mongo.UpdateInstance(filter, data),
	})
}

// FetchUserInfo returns the information about a particular user
func FetchUserInfo(c *gin.Context) {
	user := c.Param("user")
	filter := map[string]interface{}{
		"username": user,
	}
	c.JSON(200, gin.H{
		"data": mongo.FetchUserInfo(filter),
	})
}

func fetchAllInstances(instance string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		filter := map[string]interface{}{
			"instanceType": instance,
		}
		ctx.JSON(200, gin.H{
			"data": mongo.FetchDocs(mongo.InstanceCollection, filter),
		})
	}
}

// FetchAllApplications gets information for all applications deployed
func FetchAllApplications(c *gin.Context) {
	fetchAllInstances(mongo.AppInstance)(c)
}

// FetchAllDBs gets information for all applications deployed
func FetchAllDBs(c *gin.Context) {
	fetchAllInstances(mongo.DBInstance)(c)
}

// FetchAllUsers gets information for all applications deployed
func FetchAllUsers(c *gin.Context) {
	c.JSON(200, gin.H{
		"data": mongo.FetchUsers(map[string]interface{}{}),
	})
}

// DeleteUser deletes the user from database and corresponding instances
func DeleteUser(c *gin.Context) {
	user := c.Param("user")
	filter := map[string]interface{}{
		"username": user,
	}
	instanceFilter := map[string]interface{}{
		"user": user,
	}
	update := map[string]interface{}{
		"deleted": true,
	}
	go mongo.UpdateInstances(instanceFilter, update)
	c.JSON(200, gin.H{
		"message": mongo.UpdateUser(filter, update),
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
		app := c.Param("app")
		filter := map[string]interface{}{
			"name":         app,
			"instanceType": mongo.AppInstance,
		}
		update := map[string]interface{}{
			"deleted": true,
		}
		appURL, _ := redis.FetchAppNode(app)
		redis.DecrementServiceLoad(service, appURL)
		redis.RemoveApp(app)
		go commons.FullCleanup(app, mongo.AppInstance)
		c.JSON(200, gin.H{
			"message": mongo.UpdateInstance(filter, update),
		})
	}
}

// UpdateAppInfo returns a handler function for updating an application bound to a microservice
func UpdateAppInfo(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		app := c.Param("app")
		queries := c.Request.URL.Query()
		filter := utils.QueryToFilter(queries)
		if service != "dominus" {
			filter["language"] = service
		}
		filter["name"] = app
		filter["instanceType"] = mongo.AppInstance
		var data map[string]interface{}
		c.BindJSON(&data)
		c.JSON(200, gin.H{
			"message": mongo.UpdateInstance(filter, data),
		})
	}
}
