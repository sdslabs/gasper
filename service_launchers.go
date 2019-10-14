package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	"ssh": &serviceLauncher{
		Deploy: configs.ServiceConfig.SSH.Deploy,
		Start:  startSSHService("ssh"),
	},
	"ssh_proxy": &serviceLauncher{
		Deploy: configs.ServiceConfig.SSHProxy.Deploy,
		Start:  startSSHService("ssh_proxy"),
	},
	"mysql": &serviceLauncher{
		Deploy: configs.ServiceConfig.Mysql.Deploy,
		Start:  startMySQLService,
	},
	"mizu": &serviceLauncher{
		Deploy: configs.ServiceConfig.Mizu.Deploy,
		Start:  startGinService(mizu.Router, configs.ServiceConfig.Mizu.Port),
	},
	"dominus": &serviceLauncher{
		Deploy: configs.ServiceConfig.Dominus.Deploy,
		Start:  startGinService(dominus.Router, configs.ServiceConfig.Dominus.Port),
	},
	"enrai": &serviceLauncher{
		Deploy: configs.ServiceConfig.Enrai.Deploy,
		Start:  startEnraiService,
	},
	"mongodb": &serviceLauncher{
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

func startMySQLService() error {
	containers := docker.ListContainers()
	if !utils.Contains(containers, "/mysql") {
		utils.LogInfo("No Mysql instance found in host. Building the instance.")
		containerID, err := database.SetupDBInstance("mysql")
		if err != nil {
			utils.Log("There was a problem deploying MySql service.", utils.ErrorTAG)
			utils.LogError(err)
		} else {
			utils.LogInfo("Container has been deployed with ID:\t%s \n", containerID)
		}
	} else {
		containerStatus, err := docker.InspectContainerState("/mysql")
		if err != nil {
			utils.Log("Error in fetching container state. Deleting container and deploying again.", utils.ErrorTAG)
			utils.LogError(err)
			err := docker.DeleteContainer("/mysql")
			if err != nil {
				utils.LogError(err)
			}
			containerID, err := database.SetupDBInstance("mysql")
			if err != nil {
				utils.Log("There was a problem deploying MySql service even after restart.", utils.ErrorTAG)
				utils.LogError(err)
			} else {
				utils.LogInfo("Container has been deployed with ID:\t%s \n", containerID)
			}
		}
		if containerStatus["Status"].(string) == "exited" {
			err := docker.StartContainer("mysql")
			if err != nil {
				utils.LogError(err)
			}
		}
	}
	return initHTTPServer(mysql.Router, configs.ServiceConfig.Mysql.Port)
}

func startMongoDBService() error {
	containers := docker.ListContainers()
	if !utils.Contains(containers, "/mongodb") {
		utils.LogInfo("No MongoDB instance found in host. Building the instance.")
		containerID, err := database.SetupDBInstance("mongodb")
		if err != nil {
			utils.Log("There was a problem deploying mongodb service.", utils.ErrorTAG)
			utils.LogError(err)
		} else {
			utils.LogInfo("Container has been deployed with ID:\t%s \n", containerID)
		}
	} else {
		containerStatus, err := docker.InspectContainerState("/mongodb")
		if err != nil {
			utils.Log("Error in fetching container state. Deleting container and deploying again.", utils.ErrorTAG)
			utils.LogError(err)
			err := docker.DeleteContainer("/mongodb")
			if err != nil {
				utils.LogError(err)
			}
			containerID, err := database.SetupDBInstance("mongodb")
			if err != nil {
				utils.Log("There was a problem deploying MySql service even after restart.", utils.ErrorTAG)
				utils.LogError(err)
			} else {
				utils.LogInfo("Container has been deployed with ID:\t%s \n", containerID)
			}
		}
		if containerStatus["Status"].(string) == "exited" {
			err := docker.StartContainer("mongodb")
			if err != nil {
				utils.LogError(err)
			}
		}
	}
	return initHTTPServer(mongodb.Router, configs.ServiceConfig.Mongodb.Port)
}

func startSSHService(service string) func() error {
	return func() error {
		server, err := ssh.BuildSSHServer(service)
		if err != nil {
			utils.Log("There was a problem deploying SSH service. Make sure the address of Private Keys is correct in `config.json`.", utils.ErrorTAG)
			utils.LogError(err)
			return err
		}
		return server.ListenAndServe()
	}
}

func startGinService(handler *gin.Engine, port int) func() error {
	return func() error {
		return initHTTPServer(handler, port)
	}
}

func startEnraiService() error {
	return initHTTPServer(enrai.BuildEnraiServer(), configs.ServiceConfig.Enrai.Port)
}
