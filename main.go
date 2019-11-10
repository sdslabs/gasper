package main

import (
	"strings"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/enrai"
	"github.com/sdslabs/gasper/services/hikari"
	"github.com/sdslabs/gasper/services/kaze"
	"github.com/sdslabs/gasper/services/kaze/middlewares"
	"golang.org/x/sync/errgroup"
)

func initKaze() {
	go kaze.ScheduleServiceExposure()
	if configs.ServiceConfig.Kaze.Deploy {
		go kaze.ScheduleCleanup()
	}
}

func initHikari() {
	if configs.ServiceConfig.Hikari.Deploy {
		go hikari.ScheduleUpdate()
	}
}

func initEnrai() {
	if configs.ServiceConfig.Enrai.Deploy {
		go enrai.ScheduleUpdate()
	}
}

func initFalcon() {
	if configs.FalconConfig.PlugIn {
		// Initialize the Falcon Config at startup
		go middlewares.InitializeFalconConfig()
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
	initKaze()
	initHikari()
	initEnrai()
	initFalcon()
	initServices()
}
