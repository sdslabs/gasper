package database

import (
	"fmt"
	"log"

	"context"

	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	port := utils.ServiceConfig["mongoDb"].(map[string]interface{})["container_port"].(string)
	fmt.Println("che1")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%s)", port)
	fmt.Println("che2")
	ctx = context.WithValue(ctx, hostKey, agentAddress)
	ctx = context.WithValue(ctx, usernameKey, "root")
	ctx = context.WithValue(ctx, passwordKey, "root")
	//ctx = context.WithValue(ctx, databaseKey, "root")

	db, err := configDB(ctx, database, username)
	fmt.Println("che4")
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
	fmt.Println(db.Client())
	_, err = docker.ExecDetachedProcess(dbctx, cli, containerID, cmd)
	if err != nil {
		return nil
	}
	mongocmd = " mongo --eval createCollection('test') " + database
	cmd = []string{"bash", "-c", mongocmd}
	_, err = docker.ExecDetachedProcess(dbctx, cli, containerID, cmd)
	collection := db.Collection("test")
	ash := bson.D{primitive.E{Key: "autorefid", Value: "100"}}
	_, err = collection.InsertOne(ctx, ash)
	return nil
}

func configDB(ctx context.Context, database, username string) (*mongo.Database, error) {

	//url := getUri(ctx)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27018/" + username + ":" + database))
	fmt.Println("check 34")
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to mongo: %v", err)
	}
	err = client.Connect(ctx)
	fmt.Println("check 34")
	if err != nil {
		return nil, fmt.Errorf("mongo client couldn't connect with background context: %v", err)
	}
	db := client.Database(database)
	fmt.Println(db.Client())
	return db, nil
}

func DeleteMongoDB(database, username, password string) error {
	port := utils.ServiceConfig["mongoDb"].(map[string]interface{})["container_port"].(string)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%s)", port)

	ctx = context.WithValue(ctx, hostKey, agentAddress)
	ctx = context.WithValue(ctx, usernameKey, username)
	ctx = context.WithValue(ctx, passwordKey, password)
	//	ctx = context.WithValue(ctx, databaseKey, "root"

	db, err := configDB(ctx, database, username)

	err = db.Drop(ctx)
	if err != nil {
		return fmt.Errorf("Error while deleting the database : %s", err)
	}
	return nil
}
