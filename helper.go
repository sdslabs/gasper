package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/database"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/utils"
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

func buildHTTPServer(handler http.Handler, port int) *http.Server {
	if !utils.IsValidPort(port) {
		msg := fmt.Sprintf("Port %d is invalid or already in use.\n", port)
		utils.Log(msg, utils.ErrorTAG)
		panic(msg)
	}
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}

func setupDatabaseContainer(serviceName string) {
	containerName := fmt.Sprintf("/%s", serviceName)
	containers := docker.ListContainers()
	if !utils.Contains(containers, containerName) {
		utils.LogInfo("No %s instance found in host. Building the instance.", strings.Title(serviceName))
		containerID, err := database.SetupDBInstance(serviceName)
		if err != nil {
			utils.Log(fmt.Sprintf("There was a problem deploying %s service.", strings.Title(serviceName)), utils.ErrorTAG)
			utils.LogError(err)
		} else {
			utils.LogInfo("%s Container has been deployed with ID:\t%s \n", strings.Title(serviceName), containerID)
		}
	} else {
		containerStatus, err := docker.InspectContainerState(containerName)
		if err != nil {
			utils.Log("Error in fetching container state. Deleting container and deploying again.", utils.ErrorTAG)
			utils.LogError(err)
			err := docker.DeleteContainer(containerName)
			if err != nil {
				utils.LogError(err)
			}
			containerID, err := database.SetupDBInstance(serviceName)
			if err != nil {
				utils.Log(fmt.Sprintf("There was a problem deploying %s service even after restart.",
					strings.Title(serviceName)), utils.ErrorTAG)
				utils.LogError(err)
			} else {
				utils.LogInfo("Container has been deployed with ID:\t%s \n", containerID)
			}
		}
		if containerStatus["Status"].(string) == "exited" {
			err := docker.StartContainer(serviceName)
			if err != nil {
				utils.LogError(err)
			}
		}
	}
}
