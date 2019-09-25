package node

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/commons"
	g "github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func rebuildApp(c *gin.Context) {
	appName := c.Param("app")
	filter := map[string]interface{}{
		"name":         appName,
		"language":     ServiceName,
		"instanceType": "app",
	}
	data := mongo.FetchAppInfo(filter)[0]
	data["context"] = map[string]interface{}(data["context"].(primitive.D).Map())

	commons.AppFullCleanup(appName)

	resErr := pipeline(data)
	if resErr != nil {
		g.SendResponse(c, resErr, gin.H{})
		return
	}

	c.JSON(200, gin.H{
		"message": mongo.UpdateInstance(filter, data),
	})
}
