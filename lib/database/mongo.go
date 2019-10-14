package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/sdslabs/gasper/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type key string

const hostKey = key("hostKey")
const usrname = key("username")
const pwd = key("password")

var mongoUser = configs.ServiceConfig.Mongodb.Env["MONGO_INITDB_ROOT_USERNAME"].(string)
var mongoPass = configs.ServiceConfig.Mongodb.Env["MONGO_INITDB_ROOT_PASSWORD"].(string)

// CreateMongoDB creates a database in the mongodb instance with the given database name, user and password
func CreateMongoDB(database, username, password string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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
			return nil, fmt.Errorf("Error while creating the database : Database already Exists")
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
	err = (*result).Decode(v)

	if err != nil {
		errs := refreshMongoDBUser(ctx, database, username, password, client)
		if errs != nil {
			return nil, fmt.Errorf("Error while creating the database : %s", errs)
		}
	}

	return db, nil
}

// DeleteMongoDB deletes a mongo database
func DeleteMongoDB(database string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := createConnection(ctx)
	if err != nil {
		return err
	}

	db := client.Database(database)

	commandData := bson.M{"dropDatabase": 1}

	var v interface{}
	v = &(mongo.SingleResult{})
	result := db.RunCommand(ctx, commandData)
	err = (*result).Decode(v)

	if err != nil {
		return fmt.Errorf("Error while deleting the database : %s", err)
	}
	return nil
}

func createConnection(ctx context.Context) (*mongo.Client, error) {
	port := configs.ServiceConfig.Mongodb.ContainerPort
	connectionURI := fmt.Sprintf("mongodb://%s:%s@127.0.0.1:%d/admin", mongoUser, mongoPass, port)
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

func refreshMongoDBUser(ctx context.Context, database, username, password string, client *mongo.Client) error {
	db := client.Database(database)
	var v interface{}
	v = &(mongo.SingleResult{})
	dropUser := db.RunCommand(ctx, bson.M{"dropUser": username})
	err := (*dropUser).Decode(v)
	if err != nil {
		return fmt.Errorf("Error while deleting the user : %s", err)
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
