package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sdslabs/gasper/lib/database"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/seaweedfs"
	"github.com/sdslabs/gasper/lib/utils"
	"google.golang.org/grpc"
)

func checkAndPullImages(imageList ...string) {
	availableImages, err := docker.ListImages()
	if err != nil {
		utils.LogError("Main-Helper-1", err)
		os.Exit(1)
	}
	for _, image := range imageList {
		imageWithoutRepoName := strings.Replace(image, "docker.io/", "", -1)
		if utils.Contains(availableImages, image) || utils.Contains(availableImages, imageWithoutRepoName) {
			continue
		}
		utils.LogInfo("Main-Helper-2", "Image %s not present locally, pulling from DockerHUB", image)
		if err = docker.DirectPull(image); err != nil {
			utils.LogError("Main-Helper-3", err)
		}
	}
}

func startGrpcServer(server *grpc.Server, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		msg := fmt.Sprintf("Port %d is invalid or already in use", port)
		utils.Log("Main-Helper-4", msg, utils.ErrorTAG)
		os.Exit(1)
	}
	return server.Serve(lis)
}

func buildHTTPServer(handler http.Handler, port int) *http.Server {
	if !utils.IsValidPort(port) {
		msg := fmt.Sprintf("Port %d is invalid or already in use", port)
		utils.Log("Main-Helper-5", msg, utils.ErrorTAG)
		os.Exit(1)
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
	containers, err := docker.ListContainers()
	if err != nil {
		utils.LogError("Main-Helper-6", err)
		os.Exit(1)
	}

	if !utils.Contains(containers, serviceName) {
		utils.LogInfo("Main-Helper-7", "No %s instance found in host. Building the instance.", strings.Title(serviceName))
		containerID, err := database.SetupDBInstance(serviceName)
		if err != nil {
			utils.Log("Main-Helper-8", fmt.Sprintf("There was a problem deploying %s service.", strings.Title(serviceName)), utils.ErrorTAG)
			utils.LogError("Main-Helper-9", err)
		} else {
			utils.LogInfo("Main-Helper-10", "%s Container has been deployed with ID:\t%s", strings.Title(serviceName), containerID)
		}
	} else {
		containerStatus, err := docker.InspectContainerState(serviceName)
		if err != nil {
			utils.Log("Main-Helper-11", "Error in fetching container state. Deleting container and deploying again.", utils.ErrorTAG)
			utils.LogError("Main-Helper-12", err)
			err := docker.DeleteContainer(serviceName)
			if err != nil {
				utils.LogError("Main-Helper-13", err)
			}
			containerID, err := database.SetupDBInstance(serviceName)
			if err != nil {
				utils.Log("Main-Helper-14", fmt.Sprintf("There was a problem deploying %s service even after restart.",
					strings.Title(serviceName)), utils.ErrorTAG)
				utils.LogError("Main-Helper-15", err)
			} else {
				utils.LogInfo("Main-Helper-16", "Container has been deployed with ID:\t%s", containerID)
			}
		}
		if !containerStatus.Running {
			if err := docker.StartContainer(serviceName); err != nil {
				utils.LogError("Main-Helper-17", err)
			}
		}
	}
}

func setupSeaweedfsContainer(serviceName string) {
	containers, err := docker.ListContainers()
	if err != nil {
		utils.LogError("Main-Helper-18", err)
		os.Exit(1)
	}

	if !utils.Contains(containers, serviceName) {
		utils.LogInfo("No %s instance found in host. Building the instance.", strings.Title(serviceName))
		containerID, err := seaweedfs.SetupSeaweedfsInstance(serviceName)
		if err != nil {
			utils.Log("There was a problem deploying %s service.", strings.Title(serviceName), utils.ErrorTAG)
			utils.LogError("Main-Helper-19", err)
		} else {
			utils.LogInfo("%s Container has been deployed with ID:\t%s \n", strings.Title(serviceName), containerID)
		}
	} else {
		containerStatus, err := docker.InspectContainerState(serviceName)
		if err != nil {
			utils.Log("Error in fetching container state. Deleting container and deploying again.", strings.Title(serviceName), utils.ErrorTAG)
			utils.LogError("Main-Helper-20", err)
			err := docker.DeleteContainer(serviceName)
			if err != nil {
				utils.LogError("Main-Helper-21", err)
			}
			containerID, err := seaweedfs.SetupSeaweedfsInstance(serviceName)
			if err != nil {
				utils.Log("There was a problem deploying %s service even after restart.", strings.Title(serviceName), utils.ErrorTAG)
				utils.LogError("Main-Helper-22", err)
			} else {
				utils.LogInfo("Container has been deployed with ID:\t%s \n", containerID)
			}
		}
		if !containerStatus.Running {
			if err := docker.StartContainer(serviceName); err != nil {
				utils.LogError("Main-Helper-23", err)
			}
		}
	}
}
