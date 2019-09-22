package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/sdslabs/SWS/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type key string

const (
	hostKey     = key("hostKey")
	usernameKey = key("usernameKey")
	passwordKey = key("passwordKey")
)

var MongoSanitaryActionBindings = map[int]func(context.Context, string, string, string, *mongo.Client) error{
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
		log.Fatalf("database configuration failed: %v", err)
	}

	db.RunCommand(ctx , bson.M{"create":"test"})

	collection := db.Collection("test")

	ash := bson.D{primitive.E{Key: "autorefid", Value: "100"}}

	_, err = collection.InsertOne(ctx, ash)

	return nil
}

func configDB(ctx context.Context, database, username, password string) (*mongo.Database, error) {

	client, err := createConnection(ctx)
	if err != nil {
		return nil,err
	}

	dbs,errg := client.ListDatabaseNames(ctx, bson.M{})
	if errg != nil {
		return nil,fmt.Errorf("Error while creating the database : %s", errg)
	}
	for i := 0; i < len(dbs); i++ {
		if strings.Compare(dbs[i], database) == 0 {
			errs := MongoSanitaryActions(ctx, database, username, password, client, 1)
			if errs != nil {
				return nil, fmt.Errorf("Error while creating the database : %s", err)
			}
		}
	}

	db := client.Database(database)

	commandData := bson.M{ "createUser": username , 
							"pwd": password,
							"roles": bson.A {
								bson.M { "role": "dbOwner",
										 "db": database },
								"readWrite",
							},
						}
					   
	var v interface{}
	v = &(mongo.SingleResult{})
	result :=db.RunCommand(ctx,commandData)
	errh:= (*result).Decode(v)

	if errh != nil {
		errs := MongoSanitaryActions(ctx , database , username, password, client , 2)
		if errs != nil {
			return nil , fmt.Errorf("Error while creating the database : %s", err)
		}
	}

	return db, nil
}

func DeleteMongoDB(database, username string) error {
	port := configs.ServiceConfig["mongodb"].(map[string]interface{})["container_port"].(string)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%s)", port)

	ctx = context.WithValue(ctx, hostKey, agentAddress)
	ctx = context.WithValue(ctx, usernameKey, username)
	ctx = context.WithValue(ctx, passwordKey, password)

	client, err := createConnection(ctx)
	if err != nil {
		return err
	}

	db :=client.Database(database)

	err = db.Drop(ctx)
	if err != nil {
		return fmt.Errorf("Error while deleting the database : %s", err)
	}
	return nil
}

func createConnection(ctx context.Context) (*mongo.Client , error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27018/"))
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to mongo: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("mongo client couldn't connect with background context: %v", err)
	}
	return client,nil
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
	dropUser := db.RunCommand(ctx,bson.M {"dropUser" : username })
	errh:= (*dropUser).Decode(v)
		if errh != nil {
			return fmt.Errorf("Error while deleting the user : %s", errh)
		}

	commandData := bson.M{ "createUser": username , 
							"pwd": password,
							"roles": bson.A {
								bson.M { "role": "dbOwner",
										 "db": database },
								"readWrite",
							},
						}
					   
	reCreateUser:= db.RunCommand(ctx,commandData)
	errc:= (*reCreateUser).Decode(v)

	if errc != nil {
		return fmt.Errorf("Error while creating the database : %s", errc)
	}
	
	return nil
}

func MongoSanitaryActions(ctx context.Context, database, username, password string, client *mongo.Client, stage int) error {
	return MongoSanitaryActionBindings[stage](ctx, database, username, password, client)
}
