package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/database"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/utils"
	"github.com/sdslabs/SWS/services/dominus"
	"github.com/sdslabs/SWS/services/enrai"
	"github.com/sdslabs/SWS/services/mongodb"
	"github.com/sdslabs/SWS/services/mizu"
	"github.com/sdslabs/SWS/services/mysql"

	sshserver "github.com/gliderlabs/ssh"
	"github.com/sdslabs/SWS/services/ssh"
)

// UnivServer is used for handling all types of servers
type UnivServer struct {
	SSHServer  *sshserver.Server
	HTTPServer *http.Server
}

// Bind the services to the launchers
var launcherBindings = map[string]func(string, string) UnivServer{
	"ssh":     startSSHService,
	"mysql":   startMySQLService,
	"app":     startAppService,
	"enrai":   startEnraiService,
	"mongodb": startMongoDBService,
}

// Bind services to routers here
var serviceBindings = map[string]*gin.Engine{
	"dominus": dominus.Router,
	"mizu":    mizu.Router,
	"mysql":   mysql.Router,
	"mongodb": mongodb.Router,
}

func initHTTPServer(service, port string) UnivServer {
	server := &http.Server{
		Addr:         port,
		Handler:      serviceBindings[service],
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return UnivServer{
		SSHServer:  nil,
		HTTPServer: server,
	}
}

func startMySQLService(service, port string) UnivServer {
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
	server := initHTTPServer(service, port)
	return server
}

func startMongoDBService(service, port string) UnivServer {
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
	server := initHTTPServer(service, port)
	return server
}

func startSSHService(service, port string) UnivServer {
	server, err := ssh.BuildSSHServer(service)
	if err != nil {
		utils.Log("There was a problem deploying SSH service. Make sure the address of Private Keys is correct in `config.json`.", utils.ErrorTAG)
		utils.LogError(err)
		return UnivServer{
			SSHServer:  nil,
			HTTPServer: nil,
		}
	}
	return UnivServer{
		SSHServer:  server,
		HTTPServer: nil,
	}
}

func startAppService(service, port string) UnivServer {
	server := initHTTPServer(service, port)
	return server
}

func startEnraiService(service, port string) UnivServer {
	server := enrai.BuildEnraiServer(service)
	return UnivServer{
		SSHServer:  nil,
		HTTPServer: server,
	}
}

// Launcher invokes the respective launcher functions for the services
func Launcher(service, port string) UnivServer {
	if strings.HasPrefix(service, "ssh") {
		return launcherBindings["ssh"](service, port)
	} else if strings.HasPrefix(service, "enrai") {
		return launcherBindings["enrai"](service, port)
	} else if strings.HasPrefix(service, "mysql") {
		return launcherBindings["mysql"](service, port)
	} else if strings.HasPrefix(service, "mongodb") {
		return launcherBindings["mongodb"](service, port)
	} else if service != "" {
		return launcherBindings["app"](service, port)
	}

	return UnivServer{
		SSHServer:  nil,
		HTTPServer: nil,
	}
}
