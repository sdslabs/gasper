package middlewares

import (
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

var (
	errMissingCredentials   = errors.New("missing Email or Password")
	errFailedAuthentication = errors.New("incorrect Email or Password")
)

func authenticator(c *gin.Context) (interface{}, error) {
	auth := &types.Login{}
	email, _ := c.Get("Email")
	auth.Email = fmt.Sprintf("%v", email)
	user, err := mongo.FetchSingleUser(auth.GetEmail())
	if err != nil || user == nil {
		return nil, errFailedAuthentication
	}
	if !utils.CompareHashWithPassword(user.GetPassword(), auth.GetPassword()) {
		return nil, errFailedAuthentication
	}
	return user, nil
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if user, ok := data.(*types.User); ok {
		return jwt.MapClaims{
			mongo.EmailKey:    user.Email,
			mongo.UsernameKey: user.Username,
			mongo.AdminKey:    user.Admin,
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &types.User{
		Email:    claims[mongo.EmailKey].(string),
		Username: claims[mongo.UsernameKey].(string),
		Admin:    claims[mongo.AdminKey].(bool),
	}
}

func authorizator(data interface{}, c *gin.Context) bool {
	_, ok := data.(*types.User)
	return ok
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"success": false,
		"error":   message,
	})
}

// JWT handles the auth through JWT token
var JWT = &jwt.GinJWTMiddleware{
	Realm:           "Gasper",
	Key:             []byte(configs.GasperConfig.Secret),
	Timeout:         configs.JWTConfig.Timeout * time.Second,
	MaxRefresh:      configs.JWTConfig.MaxRefresh * time.Second,
	TokenLookup:     "header: Authorization",
	TokenHeadName:   "Bearer",
	TimeFunc:        time.Now,
	Authenticator:   authenticator,
	PayloadFunc:     payloadFunc,
	IdentityHandler: identityHandler,
	Authorizator:    authorizator,
	Unauthorized:    unauthorized,
}

// ExtractClaims takes the gin context and returns the User
func ExtractClaims(c *gin.Context) *types.User {
	user, success := JWT.IdentityHandler(c).(*types.User)
	if !success {
		return nil
	}
	return user
}

func init() {
	// This keeps the middleware in check if the configuration is correct
	// Prevents runtime errors
	if err := JWT.MiddlewareInit(); err != nil {
		utils.Log("Failed to initialize JWT middleware", utils.ErrorTAG)
		utils.LogError(err)
		os.Exit(1)
	}
}
