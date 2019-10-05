package middlewares

import (
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/configs"
	"github.com/sdslabs/SWS/lib/mongo"
	"golang.org/x/crypto/bcrypt"
	gojwt "gopkg.in/dgrijalva/jwt-go.v3"
)

type userStruct struct {
	Email    string
	Username string
	IsAdmin  bool
}

type authStruct struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type registerStruct struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
}

func hashPassword(password string) (string, error) {
	pass := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func compareHashWithPassword(hashedPassword, password string) bool {
	hash := []byte(hashedPassword)
	pass := []byte(password)
	err := bcrypt.CompareHashAndPassword(hash, pass)
	if err != nil {
		return false
	}
	return true
}

// Register handles registration of new users
func Register(ctx *gin.Context) {
	var register registerStruct
	if err := ctx.ShouldBind(&register); err != nil {
		ctx.JSON(400, gin.H{
			"error": "bad request parameters",
		})
		return
	}
	filter := map[string]interface{}{"email": register.Email}
	userInfo := mongo.FetchUserInfo(filter)
	if len(userInfo) > 0 {
		ctx.JSON(400, gin.H{
			"error": "email already registered",
		})
		return
	}
	hashedPass, err := hashPassword(register.Password)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err,
		})
	}
	createUser := map[string]interface{}{
		"email":    register.Email,
		"username": register.Username,
		"password": hashedPass,
		"is_admin": false,
	}
	_, err = mongo.RegisterUser(createUser)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "user created",
	})
}

// JWTMiddleware handles the auth through JWT token
var JWTMiddleware = &jwt.GinJWTMiddleware{
	Realm:         "SDS Gasper",
	Key:           []byte(configs.SWSConfig["secret"].(string)),
	Timeout:       time.Hour,
	MaxRefresh:    time.Hour,
	TokenLookup:   "header: Authorization, query: token, cookie: jwt",
	TokenHeadName: "Bearer",
	TimeFunc:      time.Now,
	Authenticator: func(ctx *gin.Context) (interface{}, error) {
		var auth authStruct
		if err := ctx.ShouldBind(&auth); err != nil {
			return nil, jwt.ErrMissingLoginValues
		}
		email := auth.Email
		password := auth.Password
		filter := map[string]interface{}{"email": email}
		userInfo := mongo.FetchUserInfo(filter)
		var userData map[string]interface{}
		if len(userInfo) == 0 {
			return nil, jwt.ErrFailedAuthentication
		}
		userData = userInfo[0]
		hashedPassword := userData["password"].(string)
		if !compareHashWithPassword(hashedPassword, password) {
			return nil, jwt.ErrFailedAuthentication
		}
		return &userStruct{
			Email:    userData["email"].(string),
			Username: userData["username"].(string),
			IsAdmin:  userData["is_admin"].(bool),
		}, nil
	},
	PayloadFunc: func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*userStruct); ok {
			return jwt.MapClaims{
				"email":    v.Email,
				"username": v.Username,
				"is_admin": v.IsAdmin,
			}
		}
		return jwt.MapClaims{}
	},
	IdentityHandler: func(mapClaims gojwt.MapClaims) interface{} {
		return &userStruct{
			Email:    mapClaims["email"].(string),
			Username: mapClaims["username"].(string),
			IsAdmin:  mapClaims["is_admin"].(bool),
		}
	},
	Authorizator: func(data interface{}, ctx *gin.Context) bool {
		_, ok := data.(*userStruct)
		return ok
	},
	Unauthorized: func(ctx *gin.Context, code int, message string) {
		ctx.JSON(code, gin.H{
			"error": message,
		})
	},
}

func init() {
	// This keeps the middleware in check if the configuration is correct
	// Prevents runtime errors
	if err := JWTMiddleware.MiddlewareInit(); err != nil {
		panic(err)
	}
}
