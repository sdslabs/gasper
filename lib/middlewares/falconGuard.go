package middlewares

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	g "github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/utils"
	falconApi "github.com/supra08/falcon-client-golang"
)

var falconConf falconApi.FalconClientGolang

<<<<<<< HEAD
// InitializeFalconConfig initializes the Falcon API with the application's credentials
=======
// InitializeFalconConfig intializes the falcon API
>>>>>>> chore(falcon): add safety check for cookie
func InitializeFalconConfig() {
	clientID := utils.FalconConfig["falconClientId"].(string)
	clientSecret := utils.FalconConfig["falconClientSecret"].(string)
	falconURLAccessToken := utils.FalconConfig["falconUrlAccessToken"].(string)
	falconURLResourceOwner := utils.FalconConfig["falconUrlResourceOwnerDetails"].(string)
	falconAccountsURL := utils.FalconConfig["falconAccountsUrl"].(string)
	falconConf = falconApi.New(clientID, clientSecret, falconURLAccessToken, falconURLResourceOwner, falconAccountsURL)
}

func getUser(cookie string) (string, error) {
	if !strings.Contains(cookie, "sdslabs") {
		return "", errors.New("User not logged in")
	}
	hash := strings.Split(cookie, "=")[1]
	user, err := falconApi.GetLoggedInUser(falconConf, hash)
	if err != nil {
		return "", errors.New("error in falcon client")
	}
	return user, nil
}

<<<<<<< HEAD
// FalconGuard returns an authorization middleware based on the plugin
=======
// FalconGuard is a middleware for checking whether the user is logged into accounts or not
>>>>>>> chore(falcon): add safety check for cookie
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
