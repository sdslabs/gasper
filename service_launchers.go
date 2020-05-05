package main

import (
	"os"
	"runtime"

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
	kaze.ServiceName: {
		Deploy: configs.ServiceConfig.Kaze.Deploy,
		Start:  startKazeService,
	},
	mizu.ServiceName: {
		Deploy: configs.ServiceConfig.Mizu.Deploy,
		Start:  startMizuService,
	},
	iwa.ServiceName: {
		Deploy: configs.ServiceConfig.Iwa.Deploy,
		Start:  startIwaService,
	},
	hikari.ServiceName: {
		Deploy: configs.ServiceConfig.Hikari.Deploy,
		Start:  hikari.NewService().ListenAndServe,
	},
	enrai.DefaultServiceName: {
		Deploy: configs.ServiceConfig.Enrai.Deploy,
		Start:  startEnraiService,
	},
	enrai.SSLServiceName: {
		Deploy: configs.ServiceConfig.Enrai.SSL.PlugIn,
		Start:  startEnraiServiceWithSSL,
	},
	kaen.ServiceName: {
		Deploy: configs.ServiceConfig.Kaen.Deploy,
		Start:  startKaenService,
	},
}

func startKaenService() error {
	if configs.ServiceConfig.Kaen.MySQL.PlugIn {
		checkAndPullImages(configs.ImageConfig.Mysql)
		setupDatabaseContainer(types.MySQL)
	}
	if configs.ServiceConfig.Kaen.MongoDB.PlugIn {
		checkAndPullImages(configs.ImageConfig.Mongodb)
		setupDatabaseContainer(types.MongoDB)
	}
	if configs.ServiceConfig.Kaen.PostgreSQL.PlugIn {
		checkAndPullImages(configs.ImageConfig.Postgresql)
		setupDatabaseContainer(types.PostgreSQL)
	}
	if configs.ServiceConfig.Kaen.RedisKaen.PlugIn {
		checkAndPullImages(configs.ImageConfig.Redis)
		setupDatabaseContainer(types.RedisKaen)
	}
	return startGrpcServer(kaen.NewService(), configs.ServiceConfig.Kaen.Port)
}

func startMizuService() error {
	images := []string{
		configs.ImageConfig.Static,
		configs.ImageConfig.Php,
		configs.ImageConfig.Nodejs,
		configs.ImageConfig.Python2,
		configs.ImageConfig.Python3,
		configs.ImageConfig.Golang,
		configs.ImageConfig.Ruby,
	}
	checkAndPullImages(images...)
	return startGrpcServer(mizu.NewService(), configs.ServiceConfig.Mizu.Port)
}

func startKazeService() error {
	if configs.ServiceConfig.Kaze.MongoDB.PlugIn {
		checkAndPullImages(configs.ImageConfig.Mongodb)
		setupDatabaseContainer(types.MongoDBGasper)
	}
	if configs.ServiceConfig.Kaze.Redis.PlugIn {
		checkAndPullImages(configs.ImageConfig.Redis)
		setupDatabaseContainer(types.RedisGasper)
	}
	return buildHTTPServer(kaze.NewService(), configs.ServiceConfig.Kaze.Port).ListenAndServe()
}

func startEnraiService() error {
	return buildHTTPServer(enrai.NewService(), configs.ServiceConfig.Enrai.Port).ListenAndServe()
}

func startIwaService() error {
	if !configs.ServiceConfig.Iwa.Deploy {
		return nil
	}
	if runtime.GOOS == "windows" {
		utils.LogInfo("Iwa doesn't work on Windows, skipping its deployment")
		return nil
	}
	return iwa.NewService().ListenAndServe()
}

func startEnraiServiceWithSSL() error {
	port := configs.ServiceConfig.Enrai.SSL.Port
	certificate := configs.ServiceConfig.Enrai.SSL.Certificate
	privateKey := configs.ServiceConfig.Enrai.SSL.PrivateKey
	err := buildHTTPServer(enrai.NewService(), port).ListenAndServeTLS(certificate, privateKey)
	if err != nil {
		utils.Log("There was a problem deploying Enrai Service with SSL", utils.ErrorTAG)
		utils.Log("Make sure the paths of certificate and private key are correct in `config.toml`", utils.ErrorTAG)
		utils.LogError(err)
		os.Exit(1)
	}
	return nil
}
