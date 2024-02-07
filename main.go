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
		go appmaker.ScheduleHealthCheck()
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

func initServices() {
	var g errgroup.Group
	for service, launcher := range launcherBindings {
		if launcher.Deploy {
			g.Go(launcher.Start)
			utils.LogInfo("Main-1", "%s Service Active", strings.Title(service))
		}
	}
	if err := g.Wait(); err != nil {
		utils.LogError("Main-2", err)
		os.Exit(1)
	}
}

func main() {
	initMaster()
	initAppMaker()
	initGenDNS()
	initGenProxy()
	initServices()
}
