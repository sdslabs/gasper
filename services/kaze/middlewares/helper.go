package middlewares

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// getBodyFromContext returns request body from a gin context without
// mutating the original body
func getBodyFromContext(c *gin.Context) []byte {
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	// Restore the io.ReadCloser to its original state
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

// Register handles registration of new users
func Register(c *gin.Context) {
	user := &types.User{}
	username, _ := c.Get("Username")
	email, _ := c.Get("Email")
	user.Email = fmt.Sprintf("%v", email)
	user.Username = fmt.Sprintf("%v", username)
	filter := types.M{mongo.EmailKey: user.Email}
	count, err := mongo.CountUsers(filter)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	if count > 0 {
		c.Next()
		return
	}
	user.SetAdmin(false)
	if _, err = mongo.RegisterUser(user); err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.Next()
}
