package main

import (
	"strings"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/dominus"
	"github.com/sdslabs/gasper/services/dominus/middlewares"
	"github.com/sdslabs/gasper/services/hikari"
	"golang.org/x/sync/errgroup"
)

func initDominus() {
	dominus.ScheduleServiceExposure()
	if configs.ServiceConfig.Dominus.Deploy {
		dominus.ScheduleCleanup()
	}
}

func initHikari() {
	if configs.ServiceConfig.Hikari.Deploy {
		hikari.ScheduleUpdate()
	}
}

func initFalcon() {
	if configs.FalconConfig.PlugIn {
		// Initialize the Falcon Config at startup
		middlewares.InitializeFalconConfig()
	}
}

func initServices() {
	var g errgroup.Group
	for service, launcher := range launcherBindings {
		if launcher.Deploy {
			g.Go(launcher.Start)
			utils.LogInfo("%s Service Active\n", strings.Title(service))
		}
	}
	if err := g.Wait(); err != nil {
		utils.LogError(err)
		panic(err)
	}
}

func main() {
	checkAndPullImages()
	initDominus()
	initHikari()
	initFalcon()
	initServices()
}
