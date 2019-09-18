package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/utils"
)

func getBodyFromContext(c *gin.Context) []byte {
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	// Restore the io.ReadCloser to its original state
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

func isUniqueInstance(instanceType, failureMessage string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data map[string]interface{}
		err := json.Unmarshal(getBodyFromContext(c), &data)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		if data["rebuild"] != nil && data["rebuild"].(bool) {
			c.Next()
			return
		}
		count, err := mongo.CountInstances(map[string]interface{}{
			"name":         data["name"].(string),
			"instanceType": instanceType,
		})
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"error": err.Error(),
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

// ValidateRequestBody validates the JSON body in a request based on the meta-data
// in the struct used to bind
func ValidateRequestBody(c *gin.Context, validationBody interface{}) {
	utils.LogDebug("Request body: %s", string(getBodyFromContext(c)))
	err := json.Unmarshal(getBodyFromContext(c), validationBody)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if result, err := validator.ValidateStruct(validationBody); !result {
		c.AbortWithStatusJSON(400, gin.H{
			"error": strings.Split(err.Error(), ";"),
		})
	} else {
		c.Next()
	}
}
