package main

import (
	"fmt"
	"net"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/dominus"
	"github.com/sdslabs/gasper/services/enrai"
	"github.com/sdslabs/gasper/services/hikari"
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
	ssh.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.SSH.Deploy,
		Start:  ssh.NewService().ListenAndServe,
	},
	mysql.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Mysql.Deploy,
		Start:  startMySQLService,
	},
	mizu.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Mizu.Deploy,
		Start:  startMizuService,
	},
	dominus.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Dominus.Deploy,
		Start:  startDominusService,
	},
	hikari.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Hikari.Deploy,
		Start:  hikari.NewService().ListenAndServe,
	},
	enrai.DefaultServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Enrai.Deploy,
		Start:  startEnraiService,
	},
	enrai.SSLServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Enrai.SSL.PlugIn,
		Start:  startEnraiServiceWithSSL,
	},
	mongodb.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Mongodb.Deploy,
		Start:  startMongoDBService,
	},
}

func startMySQLService() error {
	setupDatabaseContainer(mysql.ServiceName)
	return buildHTTPServer(mysql.NewService(), configs.ServiceConfig.Mysql.Port).ListenAndServe()
}

func startMongoDBService() error {
	setupDatabaseContainer(mongodb.ServiceName)
	return buildHTTPServer(mongodb.NewService(), configs.ServiceConfig.Mongodb.Port).ListenAndServe()
}

func startMizuService() error {
	port := configs.ServiceConfig.Mizu.Port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		msg := fmt.Sprintf("Port %d is invalid or already in use.\n", port)
		utils.Log(msg, utils.ErrorTAG)
		panic(msg)
	}
	return mizu.NewService().Serve(lis)
}

func startDominusService() error {
	return buildHTTPServer(dominus.NewService(), configs.ServiceConfig.Dominus.Port).ListenAndServe()
}

func startEnraiService() error {
	return buildHTTPServer(enrai.NewService(), configs.ServiceConfig.Enrai.Port).ListenAndServe()
}

func startEnraiServiceWithSSL() error {
	port := configs.ServiceConfig.Enrai.SSL.Port
	certificate := configs.ServiceConfig.Enrai.SSL.Certificate
	privateKey := configs.ServiceConfig.Enrai.SSL.PrivateKey
	err := buildHTTPServer(enrai.NewService(), port).ListenAndServeTLS(certificate, privateKey)
	if err != nil {
		utils.Log("There was a problem deploying Enrai Service with SSL", utils.ErrorTAG)
		utils.Log("Make sure the paths of certificate and private key are correct in `config.json`", utils.ErrorTAG)
		utils.LogError(err)
		panic(err)
	}
	return nil
}
