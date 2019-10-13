package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
)

// InsertOwner inserts the owner details into the request payload
func InsertOwner(c *gin.Context) {
	var data map[string]interface{}
	err := json.Unmarshal(getBodyFromContext(c), &data)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	data["owner"] = ExtractClaims(c).Email
	bodyBytes, err := json.Marshal(data)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Request.Header["Content-Length"] = []string{strconv.Itoa(len(bodyBytes))}
	c.Request.ContentLength = int64(len(bodyBytes))
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	c.Next()
}
