package controllers

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/kaze/middlewares"
	"github.com/sdslabs/gasper/types"
)

// Register handles registration of new users
func Register(c *gin.Context) {
	user := &types.User{}
	username, _ := c.Get("Username")
	email, _ := c.Get("Email")
	user.Email = fmt.Sprintf("%v", email)
	user.Username = fmt.Sprintf("%v", username)
	user.Password = fmt.Sprintf("%v", username)
	filter := types.M{mongo.EmailKey: user.Email}
	count, err := mongo.CountUsers(filter)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	if count > 0 {
		c.Next()
	}

	hashedPass, err := utils.HashPassword(user.GetPassword())
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	user.SetPassword(hashedPass)
	user.SetAdmin(false)

	if _, err = mongo.RegisterUser(user); err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.Next()
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
