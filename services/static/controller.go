package static

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/api"
	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/types"
)

// createApp function handles requests for making making new static app
func createApp(c *gin.Context) {
	var (
		data map[string]interface{}
	)
	c.BindJSON(&data)
	data["language"] = "static"

	conf := &types.ApplicationConfig{
		DockerImage:  "nginx:1.15.2",
		ConfFunction: configs.CreateStaticContainerConfig,
	}
	_, _ = api.CreateBasicApplication(data["name"].(string), data["github_url"].(string), "7890", "7891", conf)

	c.JSON(200, gin.H{
		"success": true,
		"id":      mongo.RegisterApp(data),
	})
}

func fetchDocs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := queryToFilter(queries)

	filter["language"] = "static"

	c.JSON(200, gin.H{
		"data": mongo.FetchAppInfo(filter),
	})
}

func deleteApp(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := queryToFilter(queries)

	filter["language"] = "static"

	c.JSON(200, gin.H{
		"message": mongo.DeleteApp(filter),
	})
}

func updateApp(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := queryToFilter(queries)

	filter["language"] = "static"

	var (
		data map[string]interface{}
	)
	c.BindJSON(&data)

	c.JSON(200, gin.H{
		"message": mongo.UpdateApp(filter, data),
	})
}

func queryToFilter(queries map[string][]string) map[string]interface{} {
	filter := make(map[string]interface{})

	for key, value := range queries {
		filter[key] = value[0]
	}

	return filter
}
