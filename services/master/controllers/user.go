package controllers

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/master/middlewares"
	"github.com/sdslabs/gasper/types"
	jwt "github.com/sdslabs/gin-jwt"
)

// Register handles registration of new users
func Register(c *gin.Context) {
	user := &types.User{}
	if err := c.BindJSON(user); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	filter := types.M{mongo.EmailKey: user.Email}
	count, err := mongo.CountUsers(filter)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	if count > 0 {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "email already registered",
		})
		return
	}

	hashedPass, err := utils.HashPassword(user.GetPassword())
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	user.SetPassword(hashedPass)
	user.SetAdmin(false)

	uuid := uuid.New()
	user.SetUUID(uuid.String())

	if _, err = mongo.RegisterUser(user); err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "user created",
	})
}

// GetUserInfo gets info regarding particular user
func GetUserInfo(c *gin.Context) {
	user, err := mongo.FetchSingleUserWithoutPassword(c.Param("user"))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(400, gin.H{
				"success": false,
				"error":   "No such user exists",
			})
			return
		}
		utils.SendServerErrorResponse(c, err)
		return
	}
	user.SetSuccess(true)
	c.JSON(200, user)
}

// GetLoggedInUserInfo returns info regarding the current logged in user
func GetLoggedInUserInfo(c *gin.Context) {
	claims := middlewares.ExtractClaims(c)
	if claims == nil {
		utils.SendServerErrorResponse(c, errors.New("Failed to extract JWT claims"))
		return
	}
	user, err := mongo.FetchSingleUserWithoutPassword(claims.GetEmail())
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	user.SetSuccess(true)
	c.JSON(200, user)
}

// UpdatePassword updates the password of a user
func UpdatePassword(c *gin.Context) {
	passwordUpdate := &types.PasswordUpdate{}
	if err := c.ShouldBind(passwordUpdate); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	claims := middlewares.ExtractClaims(c)
	if claims == nil {
		utils.SendServerErrorResponse(c, errors.New("Failed to extract JWT claims"))
		return
	}
	user, err := mongo.FetchSingleUser(claims.GetEmail())
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	if !utils.CompareHashWithPassword(user.GetPassword(), passwordUpdate.GetOldPassword()) {
		c.AbortWithStatusJSON(401, gin.H{
			"success": false,
			"error":   "old password is invalid",
		})
		return
	}
	hashedPass, err := utils.HashPassword(passwordUpdate.GetNewPassword())
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	err = mongo.UpdateUser(
		types.M{mongo.EmailKey: user.GetEmail()},
		types.M{mongo.PasswordKey: hashedPass},
	)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "password updated",
	})
}

// DeleteUser deletes the user from database
func DeleteUser(c *gin.Context) {
	claims := middlewares.ExtractClaims(c)
	if claims == nil {
		utils.SendServerErrorResponse(c, errors.New("Failed to extract JWT claims"))
		return
	}
	deleteUser(c, claims.GetEmail())
}

// GctlLogin validates the email id and alllow user to login in gctl
func GctlLogin(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	if claims == nil {
		utils.SendServerErrorResponse(c, errors.New("Failed to extract JWT claims"))
		return
	}
	exp := claims["exp"].(float64)
	tm := time.Unix(int64(exp), 0)

	email := &types.Email{}
	if err := c.Bind(email); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	filter := types.M{mongo.EmailKey: email.Email}
	count, err := mongo.CountUsers(filter)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	if count == 0 {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "email not registered",
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"expire":  tm,
	})
}

// RevokeToken updates the uuid of user so that gctl token gets invalidated
func RevokeToken(c *gin.Context) {
	auth := &types.Login{}
	if err := c.ShouldBind(auth); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	user, err := mongo.FetchSingleUser(auth.GetEmail())
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	if user == nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "email id is invalid",
		})
		return
	}
	if !utils.CompareHashWithPassword(user.GetPassword(), auth.GetPassword()) {
		c.AbortWithStatusJSON(401, gin.H{
			"success": false,
			"error":   "password is invalid",
		})
		return
	}
	uuid := uuid.New()
	err = mongo.UpdateUser(
		types.M{mongo.EmailKey: user.GetEmail()},
		types.M{mongo.GctlUUIDKey: uuid.String()},
	)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "token revoked",
	})
}
