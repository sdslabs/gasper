package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	_ "io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/factory"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/master/middlewares"
	"github.com/sdslabs/gasper/types"
)

// Endpoint to create repository in GitHub
func CreateRepository(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}

	var data *types.RepositoryRequest = &types.RepositoryRequest{}
	claims := middlewares.ExtractClaims(c)
	if claims == nil {
		utils.SendServerErrorResponse(c, errors.New("failed to extract JWT claims"))
		return
	}
	err = json.Unmarshal(raw, data)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}
	response, err := factory.CreateGithubRepository(claims.Username + data.Name)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}
	responseBody := new(bytes.Buffer)
	json.NewEncoder(responseBody).Encode(response)
	c.Data(200, "application/json", responseBody.Bytes())
}
