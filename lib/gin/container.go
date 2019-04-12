package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
)

// FetchLogs returns the container logs in a JSON format
func FetchLogs(c *gin.Context) {
	app := c.Param("app")
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	if filter["tail"] == nil {
		filter["tail"] = "-1"
	}

	appEnv, err := types.NewAppEnv()

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	data, err := docker.ReadLogs(appEnv.Context, appEnv.Client, app, filter["tail"].(string))

	if err != nil && err.Error() != "EOF" {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": data,
	})
}

// ReloadServer reloads the nginx server
func ReloadServer(c *gin.Context) {
	app := c.Param("app")
	appEnv, err := types.NewAppEnv()

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	cmd := []string{"nginx", "-s", "reload"}
	_, err = docker.ExecDetachedProcess(appEnv.Context, appEnv.Client, app, cmd)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}