package container

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/docker"
)

// Create is a controller for spawning a new container
func Create(c *gin.Context) {
	var json dockerConfig
	c.BindJSON(&json)
	docker.CreateContainer(json.Image)
}
