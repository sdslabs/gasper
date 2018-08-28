package container

import (
	"SDS/docker"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Image string `form:"image" json:"image" binding:"required"`
}

func Create(c *gin.Context) {
	var json Config
	c.BindJSON(&json)
	docker.CreateContainer(json.Image)
}
