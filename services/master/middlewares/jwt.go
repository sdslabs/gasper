package middlewares

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
	jwt "github.com/sdslabs/gin-jwt"
)

var (
	errMissingCredentials   = errors.New("missing Email or Password")
	errFailedAuthentication = errors.New("incorrect Email or Password")
)

func authenticator(c *gin.Context) (interface{}, error) {
	auth := &types.Login{}
	if err := c.ShouldBind(auth); err != nil {
		return nil, errMissingCredentials
	}
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

func payloadFuncForGctl(data interface{}) jwt.MapClaims {
	if user, ok := data.(*types.User); ok {
		return jwt.MapClaims{
			mongo.EmailKey:    user.Email,
			mongo.UsernameKey: user.Username,
			mongo.AdminKey:    user.Admin,
			mongo.GctlUUIDKey: user.GctlUUID,
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	email, ok := claims[mongo.EmailKey].(string)
	if !ok {
		return nil
	}
	username, ok := claims[mongo.UsernameKey].(string)
	if !ok {
		return nil
	}
	admin, ok := claims[mongo.AdminKey].(bool)
	if !ok {
		return nil
	}
	return &types.User{
		Email:    email,
		Username: username,
		Admin:    admin,
	}
}

func identityHandlerForGctl(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	email, ok := claims[mongo.EmailKey].(string)
	if !ok {
		return nil
	}
	username, ok := claims[mongo.UsernameKey].(string)
	if !ok {
		return nil
	}
	admin, ok := claims[mongo.AdminKey].(bool)
	if !ok {
		return nil
	}
	uuid, ok := claims[mongo.GctlUUIDKey].(string)
	if !ok {
		return nil
	}

	return &types.User{
		Email:    email,
		Username: username,
		Admin:    admin,
		GctlUUID: uuid,
	}
}

func authorizator(data interface{}, c *gin.Context) bool {
	_, ok := data.(*types.User)
	return ok
}

func authorizatorForGctl(data interface{}, c *gin.Context) bool {
	auth, ok := data.(*types.User)
	if !ok {
		return ok
	}
	user, err := mongo.FetchSingleUser(auth.GetEmail())
	if err != nil {
		return false
	}
	if auth.GctlUUID == user.GctlUUID {
		return true
	}
	return false
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

// JWTGctl handles the auth through JWT token for gctl
var JWTGctl = &jwt.GinJWTMiddleware{
	Realm:           "Gasper",
	Key:             []byte(configs.GasperConfig.Secret),
	Timeout:         configs.JWTConfig.Timeout * time.Second,
	MaxRefresh:      configs.JWTConfig.MaxRefresh * time.Second,
	TokenLookup:     "header: Authorization",
	TokenHeadName:   "gctlToken",
	TimeFunc:        time.Now,
	Authenticator:   authenticator,
	PayloadFunc:     payloadFuncForGctl,
	IdentityHandler: identityHandlerForGctl,
	Authorizator:    authorizatorForGctl,
	Unauthorized:    unauthorized,
}

//AuthRequired returns middleware according to type of request
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.Header.Get("Authorization"), "gctlToken") {
			JWTGctl.MiddlewareImpl(c)
		} else {
			JWT.MiddlewareImpl(c)
		}
	}
}

// ExtractClaims takes the gin context and returns the User
func ExtractClaims(c *gin.Context) *types.User {
	if strings.Contains(c.Request.Header.Get("Authorization"), "gctlToken") {
		user, success := JWTGctl.IdentityHandler(c).(*types.User)
		if !success {
			return nil
		}
		return user
	}

	user, success := JWT.IdentityHandler(c).(*types.User)
	if !success {
		return nil
	}
	return user
}

//LoginHandler takes the gin context and executes LoginHandler function according to authorization type
func LoginHandler(c *gin.Context) {
	if strings.Contains(c.Request.Header.Get("Authorization-Type"), "gctlToken") {
		JWTGctl.LoginHandler(c)
	} else {
		JWT.LoginHandler(c)
	}
}

//RefreshHandler takes the gin context and executes RefreshHandler function according to authorization type
func RefreshHandler(c *gin.Context) {
	if strings.Contains(c.Request.Header.Get("Authorization"), "gctlToken") {
		JWTGctl.RefreshHandler(c)
	} else {
		JWT.RefreshHandler(c)
	}
}

func init() {
	// This keeps the middleware in check if the configuration is correct
	// Prevents runtime errors
	if err := JWT.MiddlewareInit(); err != nil {
		utils.Log("Master-JWT-1", "Failed to initialize JWT middleware", utils.ErrorTAG)
		utils.LogError("Master-JWT-2", err)
		os.Exit(1)
	}
	if err := JWTGctl.MiddlewareInit(); err != nil {
		utils.Log("Master-JWT-1", "Failed to initialize JWT middleware", utils.ErrorTAG)
		utils.LogError("Master-JWT-2", err)
		os.Exit(1)
	}
}
