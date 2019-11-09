package controllers

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/middlewares"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

var immutableFields = []string{
	"name",
	"_id",
	mongo.InstanceTypeKey,
	"container_id",
	mongo.HostIPKey,
	mongo.ContainerPortKey,
	"language",
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
		"data":    mongo.FetchDocs(mongo.InstanceCollection, filter),
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

	filter["owner"] = claims.Email
	c.AbortWithStatusJSON(200, gin.H{
		"success": true,
		"data":    mongo.FetchInstances(filter),
	})
}
