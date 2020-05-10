// +build !windows

package controllers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/alphadose/gotty/backend/localcommand"
	gotty "github.com/alphadose/gotty/server"
	gottyUtils "github.com/alphadose/gotty/utils"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// DeployWebTerminal shares an application container's shell over web using `gotty`
func DeployWebTerminal(c *gin.Context) {
	appName := c.Param("app")
	port, err := utils.GetFreePort()
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	instanceURL, err := redis.FetchAppNode(appName)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Application %s is not deployed at the moment", appName),
		})
		return
	}

	if !strings.Contains(instanceURL, ":") {
		utils.SendServerErrorResponse(c,
			fmt.Errorf("Instance URL of given application is of malformed format %s", instanceURL))
		return
	}

	instanceURL = strings.Split(instanceURL, ":")[0]
	sshPort, err := redis.GetSSHPort(instanceURL)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	terminalOptions := &gotty.Options{}
	if err := gottyUtils.ApplyDefaultValues(terminalOptions); err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	terminalOptions.PermitWrite = true
	terminalOptions.Port = strconv.Itoa(port)
	terminalOptions.Once = true
	terminalOptions.Timeout = 120
	terminalOptions.TitleVariables = map[string]interface{}{
		"command":  appName,
		"argv":     "",
		"hostname": "Gasper",
	}

	backendOptions := &localcommand.Options{}
	if err := gottyUtils.ApplyDefaultValues(backendOptions); err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	termFactory, err := localcommand.NewFactory(
		"ssh", []string{"-p", sshPort, fmt.Sprintf("%s@%s", appName, instanceURL)}, backendOptions)

	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	srv, err := gotty.New(termFactory, terminalOptions)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	go srv.Run(context.Background(), gotty.WithGracefullContext(context.Background()))

	c.JSON(200, gin.H{
		"success": true,
		"url":     fmt.Sprintf("%s.%s:%d", types.Master, configs.GasperConfig.Domain, port),
		"raw_url": fmt.Sprintf("%s:%d", utils.HostIP, port),
	})
}
