package dominus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

func createApp(c *gin.Context) {
	service := c.Param("service")
	instanceURL, err := redis.GetLeastLoadedInstance(service)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	if instanceURL == "Empty Set" {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("No %s instances available at the moment", service),
		})
		return
	}
	reverseProxy(c, instanceURL)
}

func fetchDocs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)
	c.JSON(200, gin.H{
		"data": mongo.FetchAppInfo(filter),
	})
}

func execute(c *gin.Context) {
	app := c.Param("app")
	instanceURL, err := redis.FetchAppURL(app)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("Application %s is not deployed at the moment", app),
		})
		return
	}
	reverseProxy(c, instanceURL)
}
