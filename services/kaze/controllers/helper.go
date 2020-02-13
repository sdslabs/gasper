package controllers

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/kaze/middlewares"
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
}

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

// deleteUser deletes the user from database
func deleteUser(c *gin.Context, userEmail string) {
	filter := types.M{
		mongo.EmailKey: userEmail,
	}
	instanceFilter := types.M{
		mongo.OwnerKey: userEmail,
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
		"message": "user deleted",
	})
}
