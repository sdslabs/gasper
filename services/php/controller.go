package php

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/utils"
)

// phpAppConfig is json binding config for creating new php page
type phpAppConfig struct {
	Name      string `json:"name" form:"name" binding:"required"`
	UserID    int    `json:"user_id" form:"user_id" binding:"required"`
	GithubURL string `json:"github_url" form:"github_url" binding:"required"`
}

// readAndWriteConfig creates new config file for the given app
func (json phpAppConfig) ReadAndWriteConfig() utils.Error {
	// containerID, ok := os.LookupEnv("PHP_CONTAINER_ID")
	// if !ok {
	// 	return utils.Error{
	// 		Code: 500,
	// 		Err:  errors.New("PHP_CONTAINER_ID not found in the environment"),
	// 	}
	// }

	err := utils.ReadAndWriteConfig(json.Name, "php", "3b99fa7534c3")
	if err != nil {
		return utils.Error{
			Code: 500,
			Err:  err,
		}
	}

	return utils.Error{
		Code: 200,
		Err:  nil,
	}
}

// createApp function handles requests for making making new php app
func createApp(c *gin.Context) {
	var (
		json phpAppConfig
		err  utils.Error
	)
	c.BindJSON(&json)

	err = json.ReadAndWriteConfig()
	if err.Code != 200 {
		c.JSON(err.Code, gin.H{
			"message": err.Reason(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}
