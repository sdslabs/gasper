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
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/appmaker"
	"google.golang.org/grpc"
)


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
	containers := appmaker.FetchAllApplicationNames()

	if !utils.Contains(containers, serviceName) {
		utils.LogInfo("Main-Helper-6", "No %s instance found in host. Building the instance.", strings.Title(serviceName))
		containerID, err := database.SetupDBInstance(serviceName)
		if err != nil {
			utils.Log("Main-Helper-7", fmt.Sprintf("There was a problem deploying %s service.", strings.Title(serviceName)), utils.ErrorTAG)
			utils.LogError("Main-Helper-8", err)
		} else {
			utils.LogInfo("Main-Helper-9", "%s Container has been deployed with ID:\t%s", strings.Title(serviceName), containerID)
		}
	} else {
		containerStatus, err := docker.InspectContainerState(serviceName)
		if err != nil {
			utils.Log("Main-Helper-10", "Error in fetching container state. Deleting container and deploying again.", utils.ErrorTAG)
			utils.LogError("Main-Helper-11", err)
			err := docker.DeleteContainer(serviceName)
			if err != nil {
				utils.LogError("Main-Helper-12", err)
			}
			containerID, err := database.SetupDBInstance(serviceName)
			if err != nil {
				utils.Log("Main-Helper-13", fmt.Sprintf("There was a problem deploying %s service even after restart.",
					strings.Title(serviceName)), utils.ErrorTAG)
				utils.LogError("Main-Helper-14", err)
			} else {
				utils.LogInfo("Main-Helper-15", "Container has been deployed with ID:\t%s", containerID)
			}
		}
		if !containerStatus.Running {
			if err := docker.StartContainer(serviceName); err != nil {
				utils.LogError("Main-Helper-16", err)
			}
		}
	}
	// Setting up general logging for MySQL
	if serviceName == types.MySQL {
		mysqlConfig := `
[mysqld]
general_log = 1
general_log_file = /var/log/mysql/general.log
		`
		_, err := docker.ExecProcess(serviceName, []string{"sh", "-c", fmt.Sprintf("echo '%s' >> /etc/my.cnf", mysqlConfig)})
		if err != nil {
			utils.LogError("Main-Helper-17", err)
		}
		err = docker.ContainerRestart(serviceName)
		if err != nil {
			utils.LogError("Main-Helper-18", err)
		}
	}
	if serviceName == types.PostgreSQL {
		postgresConfig := `
		logging_collector = on
		log_directory = 'pg_log'
		log_filename = 'postgresql_log.log'
		log_statement = 'all'
		log_duration = on
		log_min_duration_statement = 0
		`
		_, err := docker.ExecProcess(serviceName, []string{"sh", "-c", fmt.Sprintf("echo %s >> /var/lib/postgresql/data/postgresql.conf", postgresConfig)})
		if err != nil {
			utils.LogError("Main-Helper-19", err)
		}
		err = docker.ContainerRestart(serviceName)
		if err != nil {
			utils.LogError("Main-Helper-20", err)
		}
	}
	if serviceName == types.MongoDB {
		mongoConfig := `
systemLog:
  destination: file
  logAppend: true
  path: /var/log/mongodb/mongodb.log
  verbosity: 1 `
		// command := []string{"sh", "-c", "echo '\nsystemLog:' >> /etc/mongod.conf;echo '   destination: file' >> /etc/mongod.conf;echo '   logAppend: true' >> /etc/mongod.conf;echo '   path: /var/log/mongodb/mongodb.log' >> /etc/mongod.conf;echo '   verbosity: 1' >> /etc/mongod.conf;"}
		command := []string{"sh", "-c", fmt.Sprintf("echo '%s' >> /etc/mongod.conf", mongoConfig)}
		output, err := docker.ExecProcess(serviceName, command)
		if err != nil {
			utils.LogError("Main-Helper-21 ", fmt.Errorf("Failed to update mongod.conf: %v, output: %s", err, output))
		}
		command = []string{"sh", "-c", "mongod --config /etc/mongod.conf --replSet rs0"}
		output, err = docker.ExecProcess(serviceName, command)
		if err != nil {
			utils.LogError("Main-Helper-22 ", fmt.Errorf("Failed to update mongod.conf: %v, output: %s", err, output))
		}
	}
}
