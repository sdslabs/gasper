package main

import (
	"reflect"
	"strings"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/middlewares"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/dominus"
	"golang.org/x/sync/errgroup"
)

func checkAndPullImages() {
	images := docker.ListImages()
	v := reflect.ValueOf(configs.ImageConfig)
	for i := 0; i < v.NumField(); i++ {
		image := v.Field(i).String()
		if !utils.Contains(images, image) {
			utils.LogInfo("Image %s not present locally, pulling from DockerHUB\n", image)
			docker.Pull(image)
		}
	}
}

func initDominus() {
	dominus.ScheduleServiceExposure()
	if configs.ServiceConfig.Dominus.Deploy {
		dominus.ScheduleCleanup()
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
	initFalcon()
	initServices()
}
