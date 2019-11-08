package middlewares

import (
	"encoding/json"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
	gojwt "gopkg.in/dgrijalva/jwt-go.v3"
)

const (
	emailKey    = "email"
	usernameKey = "username"
	passwordKey = "password"
	isAdminKey  = "is_admin"
)

// User to store user data after extracting from JWT Claims
type User struct {
	Email    string
	Username string
	IsAdmin  bool
}

type authBody struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type registerBody struct {
	Username string `form:"username" json:"username" binding:"required" valid:"required~Field 'username' is required but was not provided"`
	Password string `form:"password" json:"password" binding:"required" valid:"required~Field 'password' is required but was not provided"`
	Email    string `form:"email" json:"email" binding:"required" valid:"required~Field 'email' is required but was not provided,email"`
}

// RegisterValidator validates the user registration request
func RegisterValidator(ctx *gin.Context) {
	requestBody := getBodyFromContext(ctx)
	user := &registerBody{}
	err := json.Unmarshal(requestBody, user)

	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if result, err := validator.ValidateStruct(user); !result {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	ctx.Next()
}

// Register handles registration of new users
func Register(ctx *gin.Context) {
	var register registerBody
	ctx.BindJSON(&register)
	filter := types.M{emailKey: register.Email}
	userInfo := mongo.FetchUserInfo(filter)
	if len(userInfo) > 0 {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "email already registered",
		})
		return
	}
	hashedPass, err := utils.HashPassword(register.Password)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	createUser := types.M{
		emailKey:    register.Email,
		usernameKey: register.Username,
		passwordKey: hashedPass,
		isAdminKey:  false,
	}
	_, err = mongo.RegisterUser(createUser)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "User created",
		"success": true,
	})
}

// JWT handles the auth through JWT token
var JWT = &jwt.GinJWTMiddleware{
	Realm:         "SDS Gasper",
	Key:           []byte(configs.GasperConfig.Secret),
	Timeout:       time.Hour,
	MaxRefresh:    time.Hour,
	TokenLookup:   "header: Authorization",
	TokenHeadName: "Bearer",
	TimeFunc:      time.Now,
	Authenticator: func(ctx *gin.Context) (interface{}, error) {
		var auth authBody
		if err := ctx.ShouldBind(&auth); err != nil {
			return nil, jwt.ErrMissingLoginValues
		}
		email := auth.Email
		password := auth.Password
		filter := types.M{emailKey: email}
		userInfo := mongo.FetchUserInfo(filter)
		var userData types.M
		if len(userInfo) == 0 {
			return nil, jwt.ErrFailedAuthentication
		}
		userData = userInfo[0]
		hashedPassword := userData[passwordKey].(string)
		if !utils.CompareHashWithPassword(hashedPassword, password) {
			return nil, jwt.ErrFailedAuthentication
		}
		return &User{
			Email:    userData[emailKey].(string),
			Username: userData[usernameKey].(string),
			IsAdmin:  userData[isAdminKey].(bool),
		}, nil
	},
	PayloadFunc: func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*User); ok {
			return jwt.MapClaims{
				emailKey:    v.Email,
				usernameKey: v.Username,
				isAdminKey:  v.IsAdmin,
			}
		}
		return jwt.MapClaims{}
	},
	IdentityHandler: func(mapClaims gojwt.MapClaims) interface{} {
		return &User{
			Email:    mapClaims[emailKey].(string),
			Username: mapClaims[usernameKey].(string),
			IsAdmin:  mapClaims[isAdminKey].(bool),
		}
	},
	Authorizator: func(data interface{}, ctx *gin.Context) bool {
		_, ok := data.(*User)
		return ok
	},
	Unauthorized: func(ctx *gin.Context, code int, message string) {
		ctx.JSON(code, gin.H{
			"success": false,
			"error":   message,
		})
	},
}

// ExtractClaims takes the gin context and returns the User
func ExtractClaims(ctx *gin.Context) *User {
	claimsMap := jwt.ExtractClaims(ctx)
	user, success := JWT.IdentityHandler(claimsMap).(*User)
	if !success {
		return nil
	}
	return user
}

func init() {
	// This keeps the middleware in check if the configuration is correct
	// Prevents runtime errors
	if err := JWT.MiddlewareInit(); err != nil {
		panic(err)
	}
}
