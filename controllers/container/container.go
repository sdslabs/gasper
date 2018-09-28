package container

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SDS/docker"
)

// Create is a controller for spawning a new container
func Create(c *gin.Context) {
	var json createConfig
	c.BindJSON(&json)
	docker.CreateContainer(json.Image)
}
