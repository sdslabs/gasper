package main

import (
	"os"
	"strings"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/appmaker"
	"github.com/sdslabs/gasper/services/gendns"
	"github.com/sdslabs/gasper/services/genproxy"
	"github.com/sdslabs/gasper/services/master"
	"github.com/sdslabs/gasper/services/master/middlewares"
	"golang.org/x/sync/errgroup"
)

func initMaster() {
	go master.ScheduleServiceExposure()
	if configs.ServiceConfig.Master.Deploy {
		go master.ScheduleCleanup()
	}
}

func initAppMaker() {
	if configs.ServiceConfig.AppMaker.Deploy {
		go appmaker.ScheduleMetricsCollection()
	}
}

func initGenDNS() {
	if configs.ServiceConfig.GenDNS.Deploy {
		go gendns.ScheduleUpdate()
	}
}

func initGenProxy() {
	if configs.ServiceConfig.GenProxy.Deploy {
		go genproxy.ScheduleUpdate()
	}
}

func initFalcon() {
	if configs.GasperConfig.Falcon.PlugIn {
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
		os.Exit(1)
	}
}

func main() {
	initMaster()
	initAppMaker()
	initGenDNS()
	initGenProxy()
	initFalcon()
	initServices()
}
