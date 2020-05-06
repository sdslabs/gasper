package kaen

import (
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/commons"
	"github.com/sdslabs/gasper/lib/database"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/types"
)

type databaseHandler struct {
	init    func(*types.DatabaseConfig)
	create  func(types.Database) error
	delete  func(string) error
	cleanup func(string)
	logs    func(string) ([]string, error)
	reload  func() error
}

func initConstructor(language string, containerPort int) func(*types.DatabaseConfig) {
	return func(db *types.DatabaseConfig) {
		db.SetLanguage(language)
		db.SetContainerPort(containerPort)
	}
}

func cleanupConstructor(databaseType string) func(string) {
	return func(databaseName string) {
		go commons.DatabaseFullCleanup(databaseName, databaseType)
		go commons.DatabaseStateCleanup(databaseName)
	}
}

func logConstructor(serviceName string) func(string) ([]string, error) {
	return func(tail string) ([]string, error) {
		if tail == "" {
			tail = "-1"
		}
		data, err := docker.ReadLogs(serviceName, tail)
		if err != nil && err.Error() != "EOF" {
			return nil, err
		}
		return data, nil
	}
}

func reloadConstructor(serviceName string) func() error {
	return func() (err error) {
		cmd := []string{"service", serviceName, "start"}
		_, err = docker.ExecDetachedProcess(serviceName, cmd)
		return
	}
}

var pipeline = map[string]*databaseHandler{
	types.MongoDB: {
		init:    initConstructor(types.MongoDB, configs.ServiceConfig.Kaen.MongoDB.ContainerPort),
		create:  database.CreateMongoDB,
		delete:  database.DeleteMongoDB,
		cleanup: cleanupConstructor(types.MongoDB),
		logs:    logConstructor(types.MongoDB),
		reload:  reloadConstructor(types.MongoDB),
	},
	types.MySQL: {
		init:    initConstructor(types.MySQL, configs.ServiceConfig.Kaen.MySQL.ContainerPort),
		create:  database.CreateMysqlDB,
		delete:  database.DeleteMysqlDB,
		cleanup: cleanupConstructor(types.MySQL),
		logs:    logConstructor(types.MySQL),
		reload:  reloadConstructor(types.MySQL),
	},
	types.PostgreSQL: {
		init:    initConstructor(types.PostgreSQL, configs.ServiceConfig.Kaen.PostgreSQL.ContainerPort),
		create:  database.CreatePostgresqlDB,
		delete:  database.DeletePostgresqlDB,
		cleanup: cleanupConstructor(types.PostgreSQL),
		logs:    logConstructor(types.PostgreSQL),
		reload:  reloadConstructor(types.PostgreSQL),
	},
	types.Redis: {
		init:    initConstructor(types.Redis, configs.ServiceConfig.Kaen.Redis.ContainerPort),
		create:  database.CreateRedisDB,
		delete:  database.DeleteRedisDB,
		cleanup: cleanupConstructor(types.Redis),
		logs:    logConstructor(types.Redis),
		reload:  reloadConstructor(types.Redis),
	},
}
