package container

import (
	"github.com/sdslabs/SDS/docker"

	"github.com/gin-gonic/gin"
)

// Create is a controller for spawning a new container
func Create(c *gin.Context) {
	var json Config
	c.BindJSON(&json)
	docker.CreateContainer(json.Image)
}
