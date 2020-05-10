package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/utils"
)

// DeployWebTerminal shares an application container's shell over web using `gotty`
// `gotty` is not supported on windows hence an error response is returned
func DeployWebTerminal(c *gin.Context) {
	utils.SendServerErrorResponse(c, errors.New("Browser Terminal feature is not supported on Windows"))
}
