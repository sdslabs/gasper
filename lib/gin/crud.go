package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/middlewares"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
)

// FetchAppInfo returns the information about a particular app
func FetchAppInfo(c *gin.Context) {
	userStr := middlewares.ExtractClaims(c)
	app := c.Param("app")
	filter := make(map[string]interface{})
	filter["name"] = app
	filter["instanceType"] = mongo.AppInstance
	if !userStr.IsAdmin {
		filter["owner"] = userStr.Email
	}
	c.JSON(200, gin.H{
		"data": mongo.FetchAppInfo(filter),
	})
}

// FetchDBInfo returns the information about a particular db
func FetchDBInfo(c *gin.Context) {
	userStr := middlewares.ExtractClaims(c)
	db := c.Param("db")
	filter := make(map[string]interface{})
	filter["name"] = db
	filter["instanceType"] = mongo.DBInstance
	if !userStr.IsAdmin {
		filter["owner"] = userStr.Email
	}
	c.JSON(200, gin.H{
		"data": mongo.FetchDBInfo(filter),
	})
}

// UpdateAppByName updates the app getting name from url params
func UpdateAppByName(c *gin.Context) {
	userStr := middlewares.ExtractClaims(c)
	app := c.Param("app")
	filter := map[string]interface{}{
		"name":         app,
		"instanceType": mongo.AppInstance,
		"owner":        userStr.Email,
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
		"email": user,
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
		"email": user,
	}
	instanceFilter := map[string]interface{}{
		"owner": user,
	}
	update := map[string]interface{}{
		"deleted": true,
	}
	go mongo.UpdateInstances(instanceFilter, update)
	c.JSON(200, gin.H{
		"message": mongo.UpdateUser(filter, update),
	})
}

// FetchDocs fetches documents of all from mongoDB based on a filter passed
// in url query params
func FetchDocs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)
	c.JSON(200, gin.H{
		"data": mongo.FetchAppInfo(filter),
	})
}

// UpdateAppInfo updates the application document in mongoDB
func UpdateAppInfo(c *gin.Context) {
	userStr := middlewares.ExtractClaims(c)
	app := c.Param("app")
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)
	filter["name"] = app
	filter["instanceType"] = mongo.AppInstance
	filter["owner"] = userStr.Email

	var data map[string]interface{}
	c.BindJSON(&data)

	err := validateUpdatePayload(data)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": mongo.UpdateInstance(filter, data),
	})
}
