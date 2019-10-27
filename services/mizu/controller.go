package mizu

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/cloudflare"
	"github.com/sdslabs/gasper/lib/commons"
	g "github.com/sdslabs/gasper/lib/gin"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// createApp creates an application for a given language
func createApp(c *gin.Context) {
	language := c.Param("language")
	app := &types.ApplicationConfig{}
	c.BindJSON(app)

	app.DisableRebuild()
	app.SetLanguage(language)
	app.SetInstanceType(mongo.AppInstance)
	app.SetHostIP(utils.HostIP)
	app.SetNameServers(configs.GasperConfig.DNSServers)

	hikariNameServers, _ := redis.FetchServiceInstances(types.Hikari)
	for _, nameServer := range hikariNameServers {
		if strings.Contains(nameServer, ":") {
			app.AddNameServers(strings.Split(nameServer, ":")[0])
		} else {
			utils.LogError(fmt.Errorf("Hikari instance %s is of invalid format", nameServer))
		}
	}

	resErr := pipeline[language](app)
	if resErr != nil {
		g.SendResponse(c, resErr, gin.H{})
		return
	}

	sshEntrypointIP := configs.ServiceConfig.SSH.EntrypointIP
	if len(sshEntrypointIP) == 0 {
		sshEntrypointIP = utils.HostIP
	}
	app.SetSSHCmd(configs.ServiceConfig.SSH.Port, app.GetName(), sshEntrypointIP)

	if configs.CloudflareConfig.PlugIn {
		resp, err := cloudflare.CreateRecord(app.GetName(), mongo.AppInstance)
		if err != nil {
			go commons.AppFullCleanup(app.GetName())
			go commons.AppStateCleanup(app.GetName())
			utils.SendServerErrorResponse(c, err)
			return
		}
		app.SetCloudflareID(resp.Result.ID)
		app.SetAppURL(fmt.Sprintf("%s.%s.%s", app.GetName(), mongo.AppInstance, configs.GasperConfig.Domain))
	}

	err := mongo.UpsertInstance(
		types.M{
			"name":                app.GetName(),
			mongo.InstanceTypeKey: mongo.AppInstance,
		}, app)

	if err != nil && err != mongo.ErrNoDocuments {
		go commons.AppFullCleanup(app.GetName())
		go commons.AppStateCleanup(app.GetName())
		utils.SendServerErrorResponse(c, err)
		return
	}

	err = redis.RegisterApp(
		app.GetName(),
		fmt.Sprintf("%s:%d", utils.HostIP, configs.ServiceConfig.Mizu.Port),
		fmt.Sprintf("%s:%d", utils.HostIP, app.GetContainerPort()),
	)

	if err != nil {
		go commons.AppFullCleanup(app.GetName())
		go commons.AppStateCleanup(app.GetName())
		utils.SendServerErrorResponse(c, err)
		return
	}

	err = redis.IncrementServiceLoad(
		ServiceName,
		fmt.Sprintf("%s:%d", utils.HostIP, configs.ServiceConfig.Mizu.Port),
	)

	if err != nil {
		go commons.AppFullCleanup(app.GetName())
		go commons.AppStateCleanup(app.GetName())
		utils.SendServerErrorResponse(c, err)
		return
	}

	app.SetSuccess(true)
	c.JSON(200, app)
}

func rebuildApp(c *gin.Context) {
	appName := c.Param("app")

	app, err := mongo.FetchSingleApp(appName)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	commons.AppFullCleanup(appName)

	if pipeline[app.Language] == nil {
		utils.SendServerErrorResponse(c, fmt.Errorf("Non-supported language `%s` specified for `%s`", app.Language, appName))
		return
	}
	resErr := pipeline[app.Language](app)
	if resErr != nil {
		g.SendResponse(c, resErr, gin.H{})
		return
	}

	err = mongo.UpdateInstance(types.M{"name": appName}, app)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}

// deleteApp deletes an application in a worker node
func deleteApp(c *gin.Context) {
	app := c.Param("app")
	filter := types.M{
		"name":                app,
		mongo.InstanceTypeKey: mongo.AppInstance,
	}

	node, _ := redis.FetchAppNode(app)
	go redis.DecrementServiceLoad(ServiceName, node)
	go redis.RemoveApp(app)
	go commons.AppFullCleanup(app)
	if configs.CloudflareConfig.PlugIn {
		go cloudflare.DeleteRecord(app, mongo.AppInstance)
	}

	_, err := mongo.DeleteInstance(filter)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}
