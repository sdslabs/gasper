package dominus

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/redis"
)

func createApp(c *gin.Context) {
	instanceURL, err := redis.GetLeastLoadedWorker()
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if instanceURL == redis.ErrEmptySet {
		c.JSON(400, gin.H{
			"error": "No worker instances available at the moment",
		})
		return
	}
	reverseProxy(c, instanceURL)
}

func createDatabase(c *gin.Context) {
	database := c.Param("database")
	instanceURL, err := redis.GetLeastLoadedInstance(database)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if instanceURL == redis.ErrEmptySet {
		c.JSON(400, gin.H{
			"error": "No worker instances available at the moment",
		})
		return
	}
	reverseProxy(c, instanceURL)
}

func execute(c *gin.Context) {
	app := c.Param("app")
	instanceURL, err := redis.FetchAppNode(app)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("Application %s is not deployed at the moment", app),
		})
		return
	}
	reverseProxy(c, instanceURL)
}

func deleteDB(c *gin.Context) {
	db := c.Param("db")
	instanceURL, err := redis.FetchDBURL(db)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "No such database exists",
		})
		return
	}
	reverseProxy(c, instanceURL)
}
