package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/mongo"
)

// VerifyAdmin allows the request to proceed only if the user has admin privileges
func VerifyAdmin(ctx *gin.Context) {
	userStr := ExtractClaims(ctx)
	if userStr.IsAdmin {
		ctx.Next()
		return
	}
	ctx.AbortWithStatusJSON(401, gin.H{
		"error": "User does not have admin privileges",
	})
}

func isInstanceOwner(param, instanceType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		instance := ctx.Param(param)
		userStr := ExtractClaims(ctx)
		if userStr.IsAdmin {
			ctx.Next()
			return
		}

		count, err := mongo.CountInstances(map[string]interface{}{
			"name":         instance,
			"instanceType": instanceType,
			"owner":        userStr.Email,
		})
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		if count == 0 {
			ctx.AbortWithStatusJSON(401, gin.H{
				"error": fmt.Sprintf("User %s is not entitled to perform operations on %s %s", userStr.Email, param, instance),
			})
			return
		}
		ctx.Next()
	}
}

// IsAppOwner checks if a user is entitled to perform operations on an application
func IsAppOwner() gin.HandlerFunc {
	return isInstanceOwner("app", mongo.AppInstance)
}

// IsDbOwner checks if a user is entitled to perform operations on a database
func IsDbOwner() gin.HandlerFunc {
	return isInstanceOwner("db", mongo.DBInstance)
}
