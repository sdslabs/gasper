package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	g "github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/utils"
	falconApi "github.com/supra08/falcon-client-golang"
	"strings"
)

var falconConf falconApi.FalconClientGolang

func InitializeFalconConfig() {
	clientId := utils.FalconConfig["falconClientId"].(string)
	clientSecret := utils.FalconConfig["falconClientSecret"].(string)
	falconUrlAccessToken := utils.FalconConfig["falconUrlAccessToken"].(string)
	falconUrlResourceOwner := utils.FalconConfig["falconUrlResourceOwnerDetails"].(string)
	falconAccountsUrl := utils.FalconConfig["falconAccountsUrl"].(string)
	falconConf = falconApi.New(clientId, clientSecret, falconUrlAccessToken, falconUrlResourceOwner, falconAccountsUrl)
}

func getUser(cookie string) (string, error) {
	hash := strings.Split(cookie, "=")[1]
	user, err := falconApi.GetLoggedInUser(falconConf, hash)
	if err != nil {
		return "", errors.New("error in falcon client")
	}
	return user, nil
}

func FalconGuard() gin.HandlerFunc {
	if utils.FalconConfig["plugIn"].(bool) {
		return func(c *gin.Context) {
			cookie := c.GetHeader("Cookie")
			user, err := getUser(cookie)
			if user == "" {
				c.JSON(403, gin.H{
					"error": err,
				})
				c.Abort()
				return
			}
			c.Next()
		}
	}
	return func(c *g.Context) {
		c.Next()
	}
}
