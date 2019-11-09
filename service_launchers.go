package main

import (
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/enrai"
	"github.com/sdslabs/gasper/services/hikari"
	"github.com/sdslabs/gasper/services/iwa"
	"github.com/sdslabs/gasper/services/kaen"
	"github.com/sdslabs/gasper/services/kaze"
	"github.com/sdslabs/gasper/services/mizu"
	"github.com/sdslabs/gasper/types"
)

type serviceLauncher struct {
	Deploy bool
	Start  func() error
}

// Bind the services to the launchers
var launcherBindings = map[string]*serviceLauncher{
	kaze.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Kaze.Deploy,
		Start:  startKazeService,
	},
	mizu.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Mizu.Deploy,
		Start:  startMizuService,
	},
	iwa.ServiceName: &serviceLauncher{
		Deploy: configs.ServiceConfig.Iwa.Deploy,
		Start:  iwa.NewService().ListenAndServe,
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

func startKazeService() error {
	return buildHTTPServer(kaze.NewService(), configs.ServiceConfig.Kaze.Port).ListenAndServe()
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
