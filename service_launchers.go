package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/database"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/dominus"
	"github.com/sdslabs/gasper/services/enrai"
	"github.com/sdslabs/gasper/services/mizu"
	"github.com/sdslabs/gasper/services/mongodb"
	"github.com/sdslabs/gasper/services/mysql"
	"github.com/sdslabs/gasper/services/ssh"
)

type serviceLauncher struct {
	Deploy bool
	Start  func() error
}

// Bind the services to the launchers
var launcherBindings = map[string]*serviceLauncher{
	ssh.DefaultServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.SSH.Deploy,
		Start:  ssh.NewDefaultService().ListenAndServe,
	},
	ssh.ProxyServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.SSHProxy.Deploy,
		Start:  ssh.NewProxyService().ListenAndServe,
	},
	mysql.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Mysql.Deploy,
		Start:  startMySQLService,
	},
	mizu.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Mizu.Deploy,
		Start:  startHTTPService(mizu.NewService(), configs.ServiceConfig.Mizu.Port),
	},
	dominus.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Dominus.Deploy,
		Start:  startHTTPService(dominus.NewService(), configs.ServiceConfig.Dominus.Port),
	},
	enrai.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Enrai.Deploy,
		Start:  startHTTPService(enrai.NewService(), configs.ServiceConfig.Enrai.Port),
	},
	mongodb.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Mongodb.Deploy,
		Start:  startMongoDBService,
	},
}

func initHTTPServer(handler http.Handler, port int) error {
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
	return server.ListenAndServe()
}

func startHTTPService(handler http.Handler, port int) func() error {
	return func() error {
		return initHTTPServer(handler, port)
	}
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

func startMySQLService() error {
	setupDatabaseContainer(mysql.ServiceName)
	return initHTTPServer(mysql.NewService(), configs.ServiceConfig.Mysql.Port)
}

func startMongoDBService() error {
	setupDatabaseContainer(mongodb.ServiceName)
	return initHTTPServer(mongodb.NewService(), configs.ServiceConfig.Mongodb.Port)
}
