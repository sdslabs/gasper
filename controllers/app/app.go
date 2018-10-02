package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SDS/utils"
)

// CreateStaticApp function handles requests for making making new static app
func CreateStaticApp(c *gin.Context) {
	var (
		json staticAppConfig
		err  utils.Error
	)
	c.BindJSON(&json)

	err = readAndWriteStaticConf(json.Name)
	if err.Code != 200 {
		c.JSON(err.Code, gin.H{
			"message": err.Reason(),
		})
		return
	}

	c.JSON(200, gin.H{})
}
