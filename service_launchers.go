package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/services/appmaker"
	"github.com/sdslabs/gasper/services/dbmaker"
	"github.com/sdslabs/gasper/services/gendns"
	"github.com/sdslabs/gasper/services/genproxy"
	"github.com/sdslabs/gasper/services/genssh"
	"github.com/sdslabs/gasper/services/master"
	"github.com/sdslabs/gasper/types"
)

type serviceLauncher struct {
	Deploy bool
	Start  func() error
}

// Bind the services to the launchers
var launcherBindings = map[string]*serviceLauncher{
	master.ServiceName: {
		Deploy: configs.ServiceConfig.Master.Deploy,
		Start:  startMasterService,
	},
	appmaker.ServiceName: {
		Deploy: configs.ServiceConfig.AppMaker.Deploy,
		Start:  startAppMakerService,
	},
	genssh.ServiceName: {
		Deploy: configs.ServiceConfig.GenSSH.Deploy,
		Start:  startGenSSHService,
	},
	gendns.ServiceName: {
		Deploy: configs.ServiceConfig.GenDNS.Deploy,
		Start:  gendns.NewService().ListenAndServe,
	},
	genproxy.DefaultServiceName: {
		Deploy: configs.ServiceConfig.GenProxy.Deploy,
		Start:  startGenProxyService,
	},
	genproxy.SSLServiceName: {
		Deploy: configs.ServiceConfig.GenProxy.SSL.PlugIn,
		Start:  startGenProxyServiceWithSSL,
	},
	dbmaker.ServiceName: {
		Deploy: configs.ServiceConfig.DbMaker.Deploy,
		Start:  startDbMakerService,
	},
}

func startDbMakerService() error {
	if configs.ServiceConfig.DbMaker.MySQL.PlugIn {
		checkAndPullImages(configs.ImageConfig.Mysql)
		setupDatabaseContainer(types.MySQL)
	}
	if configs.ServiceConfig.DbMaker.MongoDB.PlugIn {
		checkAndPullImages(configs.ImageConfig.Mongodb)
		setupDatabaseContainer(types.MongoDB)
	}
	if configs.ServiceConfig.DbMaker.PostgreSQL.PlugIn {
		checkAndPullImages(configs.ImageConfig.Postgresql)
		setupDatabaseContainer(types.PostgreSQL)
	}
	if configs.ServiceConfig.DbMaker.Redis.PlugIn {
		checkAndPullImages(configs.ImageConfig.Redis)
	}
	return startGrpcServer(dbmaker.NewService(), configs.ServiceConfig.DbMaker.Port)
}

func startAppMakerService() error {
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
	return startGrpcServer(appmaker.NewService(), configs.ServiceConfig.AppMaker.Port)
}

func startMasterService() error {
	checkAndPullImages(configs.ImageConfig.Seaweedfs)
	err := os.MkdirAll(filepath.Join("seaweed", "seaweed-filer-storage", "filerldb2"), 0777)
	if err != nil {
		println(err.Error())
	}
	checkAndInstallSeaweedDockerPlugin()
	setupSeaweedfsContainer(types.SeaweedMaster)
	setupSeaweedfsContainer(types.SeaweedVolume)
	setupSeaweedfsContainer(types.SeaweedFiler)
	setupSeaweedfsContainer(types.SeaweedCronjob)
	setupSeaweedfsContainer(types.SeaweedS3)
	if configs.ServiceConfig.Master.MongoDB.PlugIn {
		checkAndPullImages(configs.ImageConfig.Mongodb)
		setupDatabaseContainer(types.MongoDBGasper)
	}
	if configs.ServiceConfig.Master.Redis.PlugIn {
		checkAndPullImages(configs.ImageConfig.Redis)
		setupDatabaseContainer(types.RedisGasper)
	}
	return buildHTTPServer(master.NewService(), configs.ServiceConfig.Master.Port).ListenAndServe()
}

func startGenProxyService() error {
	return buildHTTPServer(genproxy.NewService(), configs.ServiceConfig.GenProxy.Port).ListenAndServe()
}

func startGenSSHService() error {
	if !configs.ServiceConfig.GenSSH.Deploy {
		return nil
	}
	if runtime.GOOS == "windows" {
		utils.LogInfo("GenSSH doesn't work on Windows, skipping its deployment")
		return nil
	}
	return genssh.NewService().ListenAndServe()
}

func startGenProxyServiceWithSSL() error {
	port := configs.ServiceConfig.GenProxy.SSL.Port
	certificate := configs.ServiceConfig.GenProxy.SSL.Certificate
	privateKey := configs.ServiceConfig.GenProxy.SSL.PrivateKey
	err := buildHTTPServer(genproxy.NewService(), port).ListenAndServeTLS(certificate, privateKey)
	if err != nil {
		utils.Log("There was a problem deploying GenProxy Service with SSL", utils.ErrorTAG)
		utils.Log("Make sure the paths of certificate and private key are correct in `config.toml`", utils.ErrorTAG)
		utils.LogError(err)
		os.Exit(1)
	}
	return nil
}
