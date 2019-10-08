package middlewares

import "github.com/gin-gonic/gin"

// VerifyAdmin allows the request to proceed only if the user has admin privileges
func VerifyAdmin(ctx *gin.Context) {
	userStr := ExtractClaims(ctx)
	if userStr.IsAdmin {
		ctx.Next()
		return
	}
	ctx.AbortWithStatusJSON(401, gin.H{
		"error": "user does not have admin privileges",
	})
}
