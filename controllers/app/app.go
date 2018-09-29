package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SDS/utils"
)

// Create function handles requests for making making new app
func Create(c *gin.Context) {
	var response gin.H
	var err utils.Error

	appType := c.PostForm("type")

	switch appType {

	case "static":
		var json staticAppConfig
		c.BindJSON(&json)
		response, err = createStaticPage(json)

	}

	c.JSON(err.Code, response)
}
