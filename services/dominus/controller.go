package dominus

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/redis"
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
	if instanceURL == redis.ErrEmptySet {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("No %s instances available at the moment", service),
		})
		return
	}
	c.Request.URL.Path = ""
	reverseProxy(c, instanceURL)
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

func deleteDB(c *gin.Context) {
	user := c.Param("user")
	db := c.Param("db")
	dbKey := fmt.Sprintf(`%s:%s`, user, db)
	instanceURL, err := redis.FetchDBURL(dbKey)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "No such database exists",
		})
		return
	}
	reverseProxy(c, instanceURL)
}
