package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// FetchAppInfo returns the information about a particular app
func FetchAppInfo(c *gin.Context) {
	app := c.Param("app")
	filter := make(types.M)
	filter["name"] = app
	c.JSON(200, gin.H{
		"success": true,
		"data":    mongo.FetchAppInfo(filter),
	})
}

// FetchDBInfo returns the information about a particular db
func FetchDBInfo(c *gin.Context) {
	db := c.Param("db")
	filter := make(types.M)
	filter["name"] = db
	c.JSON(200, gin.H{
		"success": true,
		"data":    mongo.FetchDBInfo(filter),
	})
}

// UpdateAppByName updates the app getting name from url params
func UpdateAppByName(c *gin.Context) {
	app := c.Param("app")
	filter := types.M{
		"name":                app,
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

// FetchUserInfo returns the information about a particular user
func FetchUserInfo(c *gin.Context) {
	user := c.Param("user")
	filter := types.M{
		"email": user,
	}
	c.JSON(200, gin.H{
		"success": true,
		"data":    mongo.FetchUserInfo(filter),
	})
}

func fetchAllInstances(instance string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		filter := types.M{
			mongo.InstanceTypeKey: instance,
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
		"success": true,
		"data":    mongo.FetchUserInfo(types.M{}),
	})
}

// DeleteUser deletes the user from database and corresponding instances
func DeleteUser(c *gin.Context) {
	user := c.Param("user")
	filter := types.M{
		"email": user,
	}
	instanceFilter := types.M{
		"owner": user,
	}
	update := types.M{
		"deleted": true,
	}
	go mongo.UpdateInstances(instanceFilter, update)

	err := mongo.UpdateUser(filter, update)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

// FetchDocs fetches documents of all from mongoDB based on a filter passed
// in url query params
func FetchDocs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)
	c.JSON(200, gin.H{
		"success": true,
		"data":    mongo.FetchAppInfo(filter),
	})
}

// UpdateAppInfo updates the application document in mongoDB
func UpdateAppInfo(c *gin.Context) {
	app := c.Param("app")
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)
	filter["name"] = app
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

	err = mongo.UpdateInstance(filter, data)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}
