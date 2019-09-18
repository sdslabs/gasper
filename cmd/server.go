package main

import (
	"fmt"
	"strings"

	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/middlewares"
	"github.com/sdslabs/SWS/lib/utils"
	"github.com/sdslabs/SWS/services/dominus"
	"golang.org/x/sync/errgroup"
)

var g errgroup.Group

func main() {
	var g errgroup.Group

	images := docker.ListImages()

	for service, config := range configs.ServiceConfig {
		config := config.(map[string]interface{})
		if config["deploy"].(bool) {
			if image, check := config["image"]; check {
				image := image.(string)
				if !utils.Contains(images, image) {
					utils.Logf("Image %s not present locally, pulling from DockerHUB\n", image)
					docker.Pull(image)
				}
			}
			port := config["port"].(string)
			if utils.IsValidPort(port) {
				customServer := Launcher(service, port)
				if customServer.HTTPServer != nil {
					serviceServer := customServer.HTTPServer
					utils.Logf("%s Service Active\n", strings.Title(service))
					g.Go(func() error {
						return serviceServer.ListenAndServe()
					})
				} else if customServer.SSHServer != nil {
					serviceServer := customServer.SSHServer
					utils.Logf("%s Service Active\n", strings.Title(service))
					g.Go(func() error {
						return serviceServer.ListenAndServe()
					})
				}
			} else {
				panic(fmt.Sprintf("Cannot deploy %s service. Port %s is invalid or already in use.\n", service, port[1:]))
			}
		}
	}

	dominus.ScheduleServiceExposure()

	if configs.ServiceConfig["dominus"].(map[string]interface{})["deploy"].(bool) {
		dominus.ScheduleCleanup()
	}

	if configs.FalconConfig["plugIn"].(bool) {
		// Initialize the Falcon Config at startup
		middlewares.InitializeFalconConfig()
	}

	if err := g.Wait(); err != nil {
		panic(err)
	}
}
