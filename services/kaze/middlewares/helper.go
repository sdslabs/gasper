package middlewares

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
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
