package staticapp

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/utils"
)

// CreateApp function handles requests for making making new static app
func CreateApp(c *gin.Context) {
	var (
		json staticAppConfig
		err  utils.Error
	)
	c.BindJSON(&json)

	err = json.ReadAndWriteConfig()
	if err.Code != 200 {
		c.JSON(err.Code, gin.H{
			"message": err.Reason(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}
