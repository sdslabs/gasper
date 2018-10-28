package phpapp

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SDS/utils"
)

// CreateApp function handles requests for making making new php app
func CreateApp(c *gin.Context) {
	var (
		json phpAppConfig
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
