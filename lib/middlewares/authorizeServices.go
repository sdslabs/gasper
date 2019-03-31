package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/utils"
)

// AuthorizeService creates a gin middleware to authorize dominus requests
func AuthorizeService() gin.HandlerFunc {
	secret := utils.SWSConfig["secret"].(string)
	return func(c *gin.Context) {
		dominusSecret := c.GetHeader("dominus-secret")
		if dominusSecret == "" {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "Missing 'dominus-secret' header",
			})
			return
		}
		if dominusSecret != secret {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Invalid 'dominus-secret'",
			})
			return
		}
		c.Next()
	}
}

// ValidateParams checks whether the required parameters are present in the Request Body
func ValidateParams(paramList []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		var req map[string]interface{}
		json.Unmarshal(bodyBytes, &req)

		missingParams := []string{}
		for _, param := range paramList {
			if req[param] == nil {
				missingParams = append(missingParams, param)
			}
		}

		if len(missingParams) > 0 {
			c.AbortWithStatusJSON(400, gin.H{
				"error": fmt.Sprintf(
					"The following parameters are missing:- %s",
					strings.Join(missingParams, ",")),
			})
		} else {
			c.Next()
		}
	}
}
