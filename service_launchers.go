package main

import (
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/dominus"
	"github.com/sdslabs/gasper/services/enrai"
	"github.com/sdslabs/gasper/services/hikari"
	"github.com/sdslabs/gasper/services/kaen"
	"github.com/sdslabs/gasper/services/mizu"
	"github.com/sdslabs/gasper/services/ssh"
	"github.com/sdslabs/gasper/types"
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
	kaen.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Kaen.Deploy,
		Start:  startKaenService,
	},
}

func startKaenService() error {
	if configs.ServiceConfig.Kaen.MySQL.PlugIn {
		setupDatabaseContainer(types.MySQL)
	}
	if configs.ServiceConfig.Kaen.MongoDB.PlugIn {
		setupDatabaseContainer(types.MongoDB)
	}
	return startGrpcServer(kaen.NewService(), configs.ServiceConfig.Kaen.Port)
}

func startMizuService() error {
	return startGrpcServer(mizu.NewService(), configs.ServiceConfig.Mizu.Port)
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
