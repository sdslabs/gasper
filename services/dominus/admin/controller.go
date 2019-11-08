package admin

import (
	gogin "github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/gin"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// GetAllApplications gets all the applications from DB
func GetAllApplications(ctx *gogin.Context) {
	gin.FetchAllApplications(ctx)
}

// GetAllDatabases gets all the Databases info from DB
func GetAllDatabases(ctx *gogin.Context) {
	gin.FetchAllDBs(ctx)
}

// GetAllUsers gets all the users registered on the app
func GetAllUsers(ctx *gogin.Context) {
	gin.FetchAllUsers(ctx)
}

// GetApplicationInfo gets info regarding a particular application
func GetApplicationInfo(ctx *gogin.Context) {
	gin.FetchAppInfo(ctx)
}

// GetDatabaseInfo gets info regarding a particular database
func GetDatabaseInfo(ctx *gogin.Context) {
	gin.FetchDBInfo(ctx)
}

// GetUserInfo gets info regarding particular user
func GetUserInfo(ctx *gogin.Context) {
	gin.FetchUserInfo(ctx)
}

func changeUserPrivilege(admin bool) func(*gogin.Context) {
	return func(ctx *gogin.Context) {
		user := ctx.Param("user")
		filter := types.M{
			"email": user,
		}
		update := types.M{
			"is_admin": admin,
		}

		err := mongo.UpdateUser(filter, update)
		if err != nil {
			utils.SendServerErrorResponse(ctx, err)
			return
		}
		ctx.JSON(200, gogin.H{
			"success": true,
		})
	}
}

// GrantSuperuserPrivilege grants superuser access to a user
func GrantSuperuserPrivilege() func(ctx *gogin.Context) {
	return changeUserPrivilege(true)
}

// RevokeSuperuserPrivilege revokes superuser access from a user
func RevokeSuperuserPrivilege() func(ctx *gogin.Context) {
	return changeUserPrivilege(false)
}

// GetAllNodes fetches all the nodes registered on redis corresponding to their service
func GetAllNodes(ctx *gogin.Context) {
	services := configs.ServiceMap
	res := gogin.H{}
	// loop just to get names of services dynamically
	for service := range services {
		instances, err := redis.FetchServiceInstances(service)
		if err != nil {
			rer := types.NewResErr(500, "error when fetching from redis", err)
			gin.SendResponse(ctx, rer, gogin.H{})
			return
		}
		res[service] = instances
	}
	res["success"] = true
	ctx.JSON(200, res)
}

// GetNodesByName fetches dominus nodes for 'master' and others for 'workers'
// Rest specific service nodes are returned
func GetNodesByName(ctx *gogin.Context) {
	node := ctx.Param("type")
	res := gogin.H{}
	switch node {
	case WorkerNode:
		services := configs.ServiceMap
		for service := range services {
			if service == "dominus" {
				continue
			}
			instances, err := redis.FetchServiceInstances(service)
			if err != nil {
				rer := types.NewResErr(500, "error when fetching from redis", err)
				gin.SendResponse(ctx, rer, gogin.H{})
				return
			}
			res[service] = instances
		}
		res["success"] = true
		ctx.JSON(200, res)
		return
	case MasterNode:
		node = "dominus"
	default:
		services := configs.ServiceMap
		serviceExists := false
		for service := range services {
			if node == service {
				serviceExists = true
			}
		}
		if !serviceExists {
			rer := types.NewResErr(404, "service does not exist", nil)
			gin.SendResponse(ctx, rer, gogin.H{})
			return
		}
	}
	instances, err := redis.FetchServiceInstances(node)
	if err != nil {
		rer := types.NewResErr(500, "error when fetching from redis", err)
		gin.SendResponse(ctx, rer, gogin.H{})
		return
	}
	res[node] = instances
	res["success"] = true
	ctx.JSON(200, res)
}

// DeleteUser deletes the user from DB
func DeleteUser(ctx *gogin.Context) {
	gin.DeleteUser(ctx)
}
