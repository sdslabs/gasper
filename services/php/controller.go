package php

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/sdslabs/SWS/lib/commons"
	g "github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/mongo"
)

func rebuildApp(c *gin.Context) {
	appName := c.Param("app")
	filter := map[string]interface{}{
		"name":         appName,
		"language":     ServiceName,
		"instanceType": mongo.AppInstance,
	}
	data := mongo.FetchAppInfo(filter)[0]
	data["context"] = map[string]interface{}(data["context"].(primitive.D).Map())

	commons.FullCleanup(appName, mongo.AppInstance)

	resErr := pipeline(data)
	if resErr != nil {
		g.SendResponse(c, resErr, gin.H{})
		return
	}

	c.JSON(200, gin.H{
		"message": mongo.UpdateInstance(filter, data),
	})
}
