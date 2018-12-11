package static

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/api"
	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/types"
)

func create(c *gin.Context) {
	var json types.StaticAppConfig
	c.BindJSON(&json)

	appConf := &types.ApplicationConfig{
		DockerImage:  "nginx:1.15.2",
		ConfFunction: configs.CreateStaticContainerConfig,
	}
	err := api.CreateBasicApplication(json.Name, "7436", "7437", appConf)
	if err != nil {
		c.JSON(err.Status(), gin.H{
			"error": err.Reason(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"id":      mongo.RegisterApp(json.Name, json.UserID, json.GithubURL, "static"),
	})
}
