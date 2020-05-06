package kaen

import (
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/database"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/types"
)

// databaseStateCleanup removes the database's data from MongoDB and Redis
func databaseStateCleanup(databaseName string) {
	mongo.DeleteInstance(types.M{
		mongo.NameKey:         databaseName,
		mongo.InstanceTypeKey: mongo.DBInstance,
	})
	redis.RemoveDB(databaseName)
}

// databaseHandler is a struct for managing operations of a specific type of database (eg:- MySQL, Redis etc)
type databaseHandler struct {
	language      string
	containerPort int
	create        func(types.Database) error
	delete        func(string) error
}

// init sets the language and container port of the database server in the context
// of the new database to be created
func (handler *databaseHandler) init(db *types.DatabaseConfig) {
	db.SetLanguage(handler.language)
	db.SetContainerPort(handler.containerPort)
}

// cleanup cleans the database from MongoDB, Redis and the corresponding database server
func (handler *databaseHandler) cleanup(databaseName string) {
	go handler.delete(databaseName)
	databaseStateCleanup(databaseName)
}

// logs fetches the logs of the database server
func (handler *databaseHandler) logs(tail string) ([]string, error) {
	if tail == "" {
		tail = "-1"
	}
	data, err := docker.ReadLogs(handler.language, tail)
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}
	return data, nil
}

// reload restarts the database server
func (handler *databaseHandler) reload() error {
	cmd := []string{"service", handler.language, "start"}
	_, err := docker.ExecDetachedProcess(handler.language, cmd)
	return err
}

// pipeline maps the type of database to the corresponding handler
var pipeline = map[string]*databaseHandler{
	types.MongoDB: {
		language:      types.MongoDB,
		containerPort: configs.ServiceConfig.Kaen.MongoDB.ContainerPort,
		create:        database.CreateMongoDB,
		delete:        database.DeleteMongoDB,
	},
	types.MySQL: {
		language:      types.MySQL,
		containerPort: configs.ServiceConfig.Kaen.MySQL.ContainerPort,
		create:        database.CreateMysqlDB,
		delete:        database.DeleteMysqlDB,
	},
	types.PostgreSQL: {
		language:      types.PostgreSQL,
		containerPort: configs.ServiceConfig.Kaen.PostgreSQL.ContainerPort,
		create:        database.CreatePostgresqlDB,
		delete:        database.DeletePostgresqlDB,
	},
	types.Redis: {
		language:      types.Redis,
		containerPort: configs.ServiceConfig.Kaen.Redis.ContainerPort,
		create:        database.CreateRedisDB,
		delete:        database.DeleteRedisDB,
	},
}
