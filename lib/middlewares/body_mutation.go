package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// InsertOwner inserts the owner details into the request payload
func InsertOwner(c *gin.Context) {
	var data types.M
	err := json.Unmarshal(getBodyFromContext(c), &data)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	data["owner"] = ExtractClaims(c).Email
	bodyBytes, err := json.Marshal(data)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	c.Request.Header["Content-Length"] = []string{strconv.Itoa(len(bodyBytes))}
	c.Request.ContentLength = int64(len(bodyBytes))
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	c.Next()
}
