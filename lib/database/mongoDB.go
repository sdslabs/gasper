package database

import (
	"fmt"
	"log"

	"context"

	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type key string

const (
	hostKey     = key("hostKey")
	usernameKey = key("usernameKey")
	passwordKey = key("passwordKey")
)

// CreateMongoDB creates a database in the mongodb instance with the given database name, user and password
func CreateMongoDB(database, username, password string) error {
	port := utils.ServiceConfig["mongodb"].(map[string]interface{})["container_port"].(string)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%s)", port)
	ctx = context.WithValue(ctx, hostKey, agentAddress)
	ctx = context.WithValue(ctx, usernameKey, "root")
	ctx = context.WithValue(ctx, passwordKey, "root")
	//ctx = context.WithValue(ctx, databaseKey, "root")

	uri := getUri(ctx)
	_, err := configDB(ctx, database, uri)

	if err != nil {
		log.Fatalf("database configuration failed: %v", err)
	}

	commandData := `createUser:` + username + `, pwd: ` + password + `, roles: [
					{ role: "dbOwner", db: ` + database + ` },
					"readWrite"
			   		]`
	dbcmd := `runCommand(` + commandData + `)`

	mongocmd := fmt.Sprintf("mongo --eval %s %s", dbcmd, database)
	cmd := []string{"bash", "-c", mongocmd}

	_, err = docker.ExecDetachedProcess(dbctx, cli, containerID, cmd)
	if err != nil {
		return nil
	}

	return nil
}

func configDB(ctx context.Context, database, uri string) (*mongo.Database, error) {

	url := getUri(ctx)

	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to mongo: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("mongo client couldn't connect with background context: %v", err)
	}
	db := client.Database(database)
	return db, nil
}

func DeleteMongoDB(database, username, password string) error {
	port := utils.ServiceConfig["mongodb"].(map[string]interface{})["container_port"].(string)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%s)", port)

	ctx = context.WithValue(ctx, hostKey, agentAddress)
	ctx = context.WithValue(ctx, usernameKey, username)
	ctx = context.WithValue(ctx, passwordKey, password)
	//	ctx = context.WithValue(ctx, databaseKey, "root")
	uri := getUri(ctx)

	db, err := configDB(ctx, database, uri)
	err = db.Drop(ctx)
	if err != nil {
		return fmt.Errorf("Error while deleting the database : %s", err)
	}
	return nil
}

func getUri(ctx context.Context) string {
	uri := fmt.Sprintf(`mongodb://%s:%s@%s/%s`,
		ctx.Value(usernameKey).(string),
		ctx.Value(passwordKey).(string),
		ctx.Value(hostKey).(string),
	)
	return uri
}
