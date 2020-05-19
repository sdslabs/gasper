package middlewares

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// IsAppOwner checks if a user is entitled to perform operations on an application
func isAppOwner(c *gin.Context, instanceType string) {
	instance := c.Param("app")
	user := ExtractClaims(c)
	if user == nil {
		utils.SendServerErrorResponse(c, errors.New("Failed to extract JWT claims"))
		return
	}
	if user.IsAdmin() {
		c.Next()
		return
	}

	count, err := mongo.CountInstances(types.M{
		mongo.NameKey:         instance,
		mongo.InstanceTypeKey: instanceType,
		mongo.OwnerKey:        user.GetEmail(),
	})
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	if count == 0 {
		c.AbortWithStatusJSON(401, gin.H{
			"success": false,
			"error":   fmt.Sprintf("User %s is not entitled to perform operations on %s %s", user.GetEmail(), instanceType, instance),
		})
		return
	}
	c.Next()
}
