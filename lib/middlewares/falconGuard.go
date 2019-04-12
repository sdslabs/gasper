package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/utils"
	falconApi "github.com/supra08/falcon-client-golang"
	"strings"
)

func user(cookie string) (string, error) {

	clientId := utils.FalconConfig["falconClientId"].(string)
	clientSecret := utils.FalconConfig["falconClientSecret"].(string)
	falconUrlAccessToken := utils.FalconConfig["falconUrlAccessToken"].(string)
	falconUrlResourceOwner := utils.FalconConfig["falconUrlResourceOwnerDetails"].(string)
	falconAccountsUrl := utils.FalconConfig["falconAccountsUrl"].(string)

	config := falconApi.New(clientId, clientSecret, falconUrlAccessToken, falconUrlResourceOwner, falconAccountsUrl)
	hash := strings.Split(cookie, "=")[1]
	user, err := falconApi.GetLoggedInUser(config, hash)
	if err != nil {
		return "", errors.New("error in falcon client")
	}
	return user, nil
}

func FalconGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie := c.GetHeader("Cookie")
		user, err := user(cookie)
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
