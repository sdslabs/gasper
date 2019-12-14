package middlewares

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	g "github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	falconApi "github.com/supra08/falcon-client-golang"
)

var falconConf falconApi.FalconClientGolang

// InitializeFalconConfig intializes the falcon API
func InitializeFalconConfig() {
	clientID := configs.FalconConfig.FalconClientID
	clientSecret := configs.FalconConfig.FalconClientSecret
	falconURLAccessToken := configs.FalconConfig.FalconURLAccessToken
	falconURLResourceOwner := configs.FalconConfig.FalconURLResourceOwnerDetails
	falconAccountsURL := configs.FalconConfig.FalconAccountsURL
	falconConf = falconApi.New(clientID, clientSecret, falconURLAccessToken, falconURLResourceOwner, falconAccountsURL)
}

func getUser(cookie string, c *gin.Context) (string, error) {
	if !strings.Contains(cookie, "SDSLabs") {
		return "", errors.New("User not logged in")
	}
	hash := strings.Split(cookie, "=")[1]
	user, err := falconApi.GetLoggedInUser(falconConf, hash)
	if err != nil {
		return "", errors.New("error in falcon client")
	}
	return user, nil
}

// FalconGuard is a middleware for checking whether the user is logged into accounts or not
func FalconGuard() gin.HandlerFunc {
	if configs.FalconConfig.PlugIn {
		return func(c *gin.Context) {
			cookie := c.GetHeader("Cookie")
			user, err := getUser(cookie, c)
			if user == "" {
				c.Redirect(301, "http://arceus.sdslabs.local/")
				c.JSON(401, gin.H{
					"success": false,
					"error":   err.Error(),
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
