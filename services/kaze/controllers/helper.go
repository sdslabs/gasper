package controllers

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/kaze/middlewares"
	"github.com/sdslabs/gasper/types"
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

func fetchInstancesByUser(c *gin.Context, instanceType string) {
	filter := utils.QueryToFilter(c.Request.URL.Query())
	filter[mongo.InstanceTypeKey] = instanceType

	claims := middlewares.ExtractClaims(c)
	if claims == nil {
		utils.SendServerErrorResponse(c, errors.New("Failed to extract JWT claims"))
		return
	}

	filter[mongo.OwnerKey] = claims.Email
	c.JSON(200, gin.H{
		"success": true,
		"data":    mongo.FetchInstances(filter),
	})
}

func transferOwnership(c *gin.Context, instanceName, instanceType, newOwner string) {
	err := mongo.UpdateInstance(
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
