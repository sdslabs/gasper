package middlewares

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// VerifyAdmin allows the request to proceed only if the user has admin privileges
func VerifyAdmin(ctx *gin.Context) {
	claims := ExtractClaims(ctx)
	if claims == nil {
		utils.SendServerErrorResponse(ctx, errors.New("Failed to extract JWT claims"))
		return
	}
	if claims.IsAdmin {
		ctx.Next()
		return
	}
	ctx.AbortWithStatusJSON(401, gin.H{
		"success": false,
		"error":   "User does not have admin privileges",
	})
}

func isInstanceOwner(instanceType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		instance := ctx.Param(instanceType)
		claims := ExtractClaims(ctx)
		if claims == nil {
			utils.SendServerErrorResponse(ctx, errors.New("Failed to extract JWT claims"))
			return
		}
		if claims.IsAdmin {
			ctx.Next()
			return
		}
		count, err := mongo.CountInstances(types.M{
			"name":                instance,
			mongo.InstanceTypeKey: instanceType,
			"owner":               claims.Email,
		})
		if err != nil {
			utils.SendServerErrorResponse(ctx, err)
			return
		}

		if count == 0 {
			ctx.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"error":   fmt.Sprintf("User %s is not entitled to perform operations on %s %s", claims.Email, instanceType, instance),
			})
			return
		}
		ctx.Next()
	}
}

// IsAppOwner checks if a user is entitled to perform operations on an application
func IsAppOwner() gin.HandlerFunc {
	return isInstanceOwner(mongo.AppInstance)
}

// IsDbOwner checks if a user is entitled to perform operations on a database
func IsDbOwner() gin.HandlerFunc {
	return isInstanceOwner(mongo.DBInstance)
}
