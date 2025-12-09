package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/master/middlewares"
	"github.com/sdslabs/gasper/types"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var immutableFields = []string{
	mongo.NameKey,
	"_id",
	mongo.InstanceTypeKey,
	"container_id",
	mongo.HostIPKey,
	mongo.ContainerPortKey,
	mongo.LanguageKey,
	"cloudflare_id",
	"app_url",
	"docker_image",
	// "repo_url",
}

// ValidateUpdatePayload was used to validate the update payload
// Deprecated : This is no longer used as the update payload is now parsed into a struct
// Instead use middleware.ValidateApplicationUpdateRequest instead
func validateUpdatePayload(data types.M) error {
	res := ""
	for _, field := range immutableFields {
		if data[field] != nil {
			res += fmt.Sprintf("Field `%s` is immutable; ", field)
		}
	}

	if res != "" {
		return errors.New(res)
	}
	return nil
}

// UpdateData updates the data of an application using the update request payload
func UpdateData(app *types.ApplicationConfig, data *types.UpdatePayload) error {
	if data.Password != nil {
		app.Password = *data.Password
	}
	if data.Git != nil {
		if data.Git.AccessToken != nil {
			app.Git.AccessToken = *data.Git.AccessToken
		}
		if data.Git.Branch != nil {
			app.Git.Branch = *data.Git.Branch
		}
	}
	if data.Context != nil {
		if data.Context.Index != nil {
			app.Context.Index = *data.Context.Index
		}
		if data.Context.RcFile != nil {
			app.Context.RcFile = *data.Context.RcFile
		}
		if data.Context.Build != nil {
			app.Context.Build = *data.Context.Build
		}
		if data.Context.Run != nil {
			app.Context.Run = *data.Context.Run
		}
	}

	if data.Resources != nil {
		ok, err := utils.ValidateCPUvalue(data.Resources.CPU)
		if ok && err == nil {
			app.Resources.CPU = data.Resources.CPU
		} else {
			return err
		}
		ok, err = utils.ValidateRAMvalue(data.Resources.Memory)
		if ok && err == nil {
			app.Resources.Memory = data.Resources.Memory
		} else {
			return err
		}
	}
	return nil
}

func fetchInstances(c *gin.Context, instance string) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)
	filter[mongo.InstanceTypeKey] = instance
	c.JSON(200, gin.H{
		"success": true,
		"data":    mongo.FetchDocs(mongo.InstanceCollection, filter, nil),
	})
}

// FetchAllInstancesByUser returns all instances owned by a user
func FetchAllInstancesByUser(c *gin.Context) {
	filter := utils.QueryToFilter(c.Request.URL.Query())
	claims := middlewares.ExtractClaims(c)
	if claims == nil {
		utils.SendServerErrorResponse(c, errors.New("Failed to extract JWT claims"))
		return
	}
	filter[mongo.OwnerKey] = claims.GetEmail()

	projection := types.M{
		mongo.NameKey:         1,
		mongo.LanguageKey:     1,
		mongo.InstanceTypeKey: 1,
	}

	opts := options.Find().SetProjection(projection)

	c.JSON(200, gin.H{
		"success": true,
		"data":    mongo.FetchDocs(mongo.InstanceCollection, filter, opts),
	})
}

func fetchInstancesByUser(c *gin.Context, instanceType string) {
	filter := utils.QueryToFilter(c.Request.URL.Query())
	filter[mongo.InstanceTypeKey] = instanceType

	claims := middlewares.ExtractClaims(c)
	if claims == nil {
		utils.SendServerErrorResponse(c, errors.New("Failed to extract JWT claims"))
		return
	}

	filter[mongo.OwnerKey] = claims.GetEmail()
	c.JSON(200, gin.H{
		"success": true,
		"data":    mongo.FetchInstances(filter),
	})
}

func transferOwnership(c *gin.Context, instanceName, instanceType, newOwner string) {
	count, err := mongo.CountUsers(types.M{mongo.EmailKey: newOwner})
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	if count == 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Recipent's email address is invalid",
		})
		return
	}
	err = mongo.UpdateInstance(
		types.M{
			mongo.NameKey:         instanceName,
			mongo.InstanceTypeKey: instanceType,
		},
		types.M{
			mongo.OwnerKey: newOwner,
		},
	)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

// Handle404 handles 404 errors
func Handle404(c *gin.Context) {
	c.AbortWithStatusJSON(404, gin.H{
		"success": false,
		"error":   "Page not found",
	})
}

func deleteinstance(c *gin.Context, appName string, wg *sync.WaitGroup) {
	defer wg.Done()
	DeleteAppUsingAppname(c, appName)
}

// deleteUser deletes the user from database
func deleteUser(c *gin.Context, userEmail string) {
	filter := types.M{
		mongo.EmailKey: userEmail,
	}
	instanceFilter := types.M{
		mongo.OwnerKey: userEmail,
	}

	instancesInfo := mongo.FetchInstances(instanceFilter)
	var wg sync.WaitGroup
	for _, instanceData := range instancesInfo {
		appName := instanceData["name"].(string)
		wg.Add(1)
		go deleteinstance(c, appName, &wg)
	}

	wg.Wait()
	go mongo.DeleteInstances(instanceFilter)

	_, err := mongo.DeleteUser(filter)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "user deleted",
	})
}

func appsOnNode(hostIP string) ([]types.ApplicationConfig, error) {
	apps := mongo.FetchAppInfo(types.M{mongo.HostIPKey: hostIP})
	var appObjects []types.ApplicationConfig
	for _, app := range apps {
		var appObject types.ApplicationConfig
		temp, err := json.Marshal(app)
		if err != nil {
			utils.LogError("Marshalling Error", err)
			return nil, err
		}
		err = json.Unmarshal(temp, &appObject)
		if err != nil {
			utils.LogError("Unmarshalling Error", err)
			return nil, err
		}
		appObjects = append(appObjects, appObject)
	}
	return appObjects, nil
}
