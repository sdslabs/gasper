package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	falconApi "github.com/supra08/falcon-client-golang"
	"strings"
)

func user(cookie string) (string, error) {
	config := falconApi.New("howl-MKUlTqXmtQHPEtPN", "0b3d4a96a621bf4c1d08ceddc49241482c96aa1b7e69c3e81f5f4190c80c0d8b", "http://falcon.sdslabs.local/access_token", "http://falcon.sdslabs.local/users/", "http://arceus.sdslabs.local/")
	hash := strings.Split(cookie, "=")[1]
	fmt.Println(hash)
	user, err := falconApi.GetLoggedInUser(config, hash)
	if err != nil {
		return "", fmt.Errorf("error in falcon client")
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
