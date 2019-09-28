package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/sdslabs/SWS/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type key string

const hostKey = key("hostKey")

var mongoUser = configs.ServiceConfig["mongodb"].(map[string]interface{})["env"].(map[string]interface{})["MONGO_INITDB_ROOT_USERNAME"].(string)
var mongoPass = configs.ServiceConfig["mongodb"].(map[string]interface{})["env"].(map[string]interface{})["MONGO_INITDB_ROOT_PASSWORD"].(string)

var mongoSanitaryActionBindings = map[int]func(context.Context, string, string, string, *mongo.Client) error{
	1: refreshMongoDB,
	2: refreshMongoDBUser,
}

// CreateMongoDB creates a database in the mongodb instance with the given database name, user and password
func CreateMongoDB(database, username, password string) error {
	port := configs.ServiceConfig["mongodb"].(map[string]interface{})["container_port"].(string)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%s)", port)
	ctx = context.WithValue(ctx, hostKey, agentAddress)

	_, err := configDB(ctx, database, username, password)

	if err != nil {
		return fmt.Errorf("database configuration failed: %v", err)
	}

	return nil
}

func configDB(ctx context.Context, database, username, password string) (*mongo.Database, error) {

	client, err := createConnection(ctx)
	if err != nil {
		return nil, err
	}

	dbs, errg := client.ListDatabaseNames(ctx, bson.M{})

	if errg != nil {
		return nil, fmt.Errorf("Error while creating the database : %s", errg)
	}

	for i := 0; i < len(dbs); i++ {
		if strings.Compare(dbs[i], database) == 0 {
			errs := mongoSanitaryActions(ctx, database, username, password, client, 1)
			if errs != nil {
				return nil, fmt.Errorf("Error while creating the database : %s", err)
			}
		}
	}

	db := client.Database(database)

	commandData := bson.M{"createUser": username,
		"pwd": password,
		"roles": bson.A{
			bson.M{"role": "dbOwner",
				"db": database},
			"readWrite",
		},
	}

	var v interface{}
	v = &(mongo.SingleResult{})
	result := db.RunCommand(ctx, commandData)
	errh := (*result).Decode(v)

	if errh != nil {
		errs := mongoSanitaryActions(ctx, database, username, password, client, 2)
		if errs != nil {
			return nil, fmt.Errorf("Error while creating the database : %s", err)
		}
	}

	return db, nil
}

// DeleteMongoDB deletes a mongo database
func DeleteMongoDB(database, username string) error {
	port := configs.ServiceConfig["mongodb"].(map[string]interface{})["container_port"].(string)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%s)", port)

	ctx = context.WithValue(ctx, hostKey, agentAddress)

	client, err := createConnection(ctx)
	if err != nil {
		return err
	}

	db := client.Database(database)

	err = db.Drop(ctx)
	if err != nil {
		return fmt.Errorf("Error while deleting the database : %s", err)
	}
	return nil
}

func createConnection(ctx context.Context) (*mongo.Client, error) {
	port := configs.ServiceConfig["mongodb"].(map[string]interface{})["container_port"].(string)
	connectionURI := fmt.Sprintf("mongodb://%s:%s@127.0.0.1:%s/", mongoUser, mongoPass, port)
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to mongo: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("mongo client couldn't connect with background context: %v", err)
	}
	return client, nil
}

func refreshMongoDB(ctx context.Context, database, username, password string, client *mongo.Client) error {
	db := client.Database(database)
	err := db.Drop(ctx)
	if err != nil {
		return fmt.Errorf("Error while deleting the database : %s", err)
	}
	return nil
}

func refreshMongoDBUser(ctx context.Context, database, username, password string, client *mongo.Client) error {
	db := client.Database(database)
	var v interface{}
	v = &(mongo.SingleResult{})
	dropUser := db.RunCommand(ctx, bson.M{"dropUser": username})
	errh := (*dropUser).Decode(v)
	if errh != nil {
		return fmt.Errorf("Error while deleting the user : %s", errh)
	}

	commandData := bson.M{"createUser": username,
		"pwd": password,
		"roles": bson.A{
			bson.M{"role": "dbOwner",
				"db": database},
			"readWrite",
		},
	}

	reCreateUser := db.RunCommand(ctx, commandData)
	errc := (*reCreateUser).Decode(v)

	if errc != nil {
		return fmt.Errorf("Error while creating the database : %s", errc)
	}

	return nil
}

func mongoSanitaryActions(ctx context.Context, database, username, password string, client *mongo.Client, stage int) error {
	return mongoSanitaryActionBindings[stage](ctx, database, username, password, client)
}
