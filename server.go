package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/middlewares"
	"github.com/sdslabs/SWS/lib/utils"
	"github.com/sdslabs/SWS/services/dominus"
	"github.com/sdslabs/SWS/services/node"
	"github.com/sdslabs/SWS/services/php"
	"github.com/sdslabs/SWS/services/python"
	"github.com/sdslabs/SWS/services/ssh"
	"github.com/sdslabs/SWS/services/static"
	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group

	// Bind services to routers here
	serviceBindings := map[string]*gin.Engine{
		"dominus": dominus.Router,
		"static":  static.Router,
		"php":     php.Router,
		"node":    node.Router,
		"python":  python.Router,
	}

	images := docker.ListImages()

	for service, config := range utils.ServiceConfig {
		config := config.(map[string]interface{})
		if config["deploy"].(bool) {
			if image, check := config["image"]; check {
				image := image.(string)
				if !utils.Contains(images, image) {
					fmt.Printf("Image %s not present locally, pulling from DockerHUB\n", image)
					docker.Pull(image)
				}
			}
			port := config["port"].(string)
			if utils.IsValidPort(port) {
				if strings.HasPrefix(service, "ssh") {
					server, err := ssh.BuildSSHServer(service)
					if err != nil {
						fmt.Println("There was a problem deploying SSH service. Make sure the address of Private Keys is correct in `config.json`.")
						fmt.Printf("ERROR:: %s\n", err.Error())
					} else {
						fmt.Printf("%s Service Active\n", strings.Title(service))
						g.Go(func() error {
							return server.ListenAndServe()
						})
					}
				} else {
					server := &http.Server{
						Addr:         config["port"].(string),
						Handler:      serviceBindings[service],
						ReadTimeout:  5 * time.Second,
						WriteTimeout: 30 * time.Second,
					}
					fmt.Printf("%s Service Active\n", strings.Title(service))
					g.Go(func() error {
						return server.ListenAndServe()
					})
				}
			} else {
				panic(fmt.Sprintf("Cannot deploy %s service. Port %s is invalid or already in use.\n", service, port[1:]))
			}
		}
	}

	dominus.ExposeServices()

	if utils.ServiceConfig["dominus"].(map[string]interface{})["deploy"].(bool) {
		cleanupInterval := time.Duration(utils.SWSConfig["cleanupInterval"].(float64))
		dominus.ScheduleCleanup(cleanupInterval * time.Second)
	}

	if err := g.Wait(); err != nil {
		panic(err)
	}

	if utils.FalconConfig["plugIn"].(bool) {
		// Initialize the Falcon Config at startup
		middlewares.InitializeFalconConfig()
	}
}
