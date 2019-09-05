package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/mongo"
)

func isUniqueInstance(instanceType, failureMessage string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		var data map[string]interface{}
		err := json.Unmarshal(bodyBytes, &data)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "Invalid JSON",
			})
			return
		}
		count, err := mongo.CountInstances(map[string]interface{}{
			"name":         data["name"].(string),
			"instanceType": instanceType,
		})
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"error": err,
			})
			return
		}
		if count != 0 {
			c.AbortWithStatusJSON(400, gin.H{
				"error": failureMessage,
			})
			return
		}
		c.Next()
	}
}

// IsUniqueApp checks whether the application name is unique or not
func IsUniqueApp() gin.HandlerFunc {
	return isUniqueInstance(mongo.AppInstance, "Application with that name already exists")
}

// IsUniqueDB checks whether the database name is unique or not
func IsUniqueDB() gin.HandlerFunc {
	return isUniqueInstance(mongo.DBInstance, "Database with that name already exists")
}
