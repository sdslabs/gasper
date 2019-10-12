package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/types"
	"github.com/sdslabs/gasper/lib/utils"
)

// SendResponse sends the response back to client
// r is nil when status code is OK (200)
// response should be set when r is not nil
func SendResponse(c *gin.Context, r types.ResponseError, response gin.H) {
	if r == nil { // OK
		c.JSON(200, response)
		return
	}
	utils.LogResErr(r)
	if r.Status() == 500 {
		if configs.GasperConfig.Debug {
			c.JSON(500, gin.H{
				"error": r.Verbose(),
			})
		} else {
			c.JSON(500, gin.H{
				"error": "INTERNAL_SERVER_ERROR",
			})
		}
		return
	}
	c.JSON(r.Status(), gin.H{
		"error": r.Message(),
	})
}
