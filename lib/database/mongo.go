package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoRootUser     = configs.ServiceConfig.Mongodb.Env["MONGO_INITDB_ROOT_USERNAME"].(string)
	mongoRootPassword = configs.ServiceConfig.Mongodb.Env["MONGO_INITDB_ROOT_PASSWORD"].(string)
)

// CreateMongoDB creates a database in the mongodb instance with the given database name, user and password
func CreateMongoDB(db types.Database) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	_, err := configDB(ctx, db)

	if err != nil {
		return fmt.Errorf("database configuration failed: %v", err)
	}

	return nil
}

func configDB(ctx context.Context, db types.Database) (*mongo.Database, error) {
	client, err := createConnection(ctx)
	if err != nil {
		return nil, err
	}

	dbs, err := client.ListDatabaseNames(ctx, bson.M{})

	if err != nil {
		return nil, fmt.Errorf("Error while creating the database : %s", err)
	}

	for i := 0; i < len(dbs); i++ {
		if strings.Compare(dbs[i], db.GetName()) == 0 {
			return nil, fmt.Errorf("Error while creating the database : Database already Exists")
		}
	}

	conn := client.Database(db.GetName())

	commandData := bson.M{"createUser": db.GetUser(),
		"pwd": db.GetPassword(),
		"roles": bson.A{
			bson.M{"role": "dbOwner",
				"db": db.GetName()},
			"readWrite",
		},
	}

	var v interface{}
	v = &(mongo.SingleResult{})
	result := conn.RunCommand(ctx, commandData)

	if err = (*result).Decode(v); err != nil {
		if err = refreshMongoDBUser(ctx, db, client); err != nil {
			return nil, fmt.Errorf("Error while creating the database : %s", err)
		}
	}

	return conn, nil
}

// DeleteMongoDB deletes a mongo database
func DeleteMongoDB(databaseName string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := createConnection(ctx)
	if err != nil {
		return err
	}

	conn := client.Database(databaseName)

	commandData := bson.M{"dropDatabase": 1}

	var v interface{}
	v = &(mongo.SingleResult{})
	result := conn.RunCommand(ctx, commandData)
	err = (*result).Decode(v)

	if err != nil {
		return fmt.Errorf("Error while deleting the database : %s", err)
	}
	return nil
}

func createConnection(ctx context.Context) (*mongo.Client, error) {
	port := configs.ServiceConfig.Mongodb.ContainerPort
	connectionURI := fmt.Sprintf("mongodb://%s:%s@127.0.0.1:%d/admin", mongoRootUser, mongoRootPassword, port)
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

func refreshMongoDBUser(ctx context.Context, db types.Database, client *mongo.Client) error {
	conn := client.Database(db.GetName())
	var v interface{}
	v = &(mongo.SingleResult{})
	dropUser := conn.RunCommand(ctx, bson.M{"dropUser": db.GetUser()})
	err := (*dropUser).Decode(v)
	if err != nil {
		return fmt.Errorf("Error while deleting the user : %s", err)
	}

	commandData := bson.M{"createUser": db.GetUser(),
		"pwd": db.GetPassword(),
		"roles": bson.A{
			bson.M{"role": "dbOwner",
				"db": db.GetName()},
			"readWrite",
		},
	}

	reCreateUser := conn.RunCommand(ctx, commandData)
	err = (*reCreateUser).Decode(v)

	if err != nil {
		return fmt.Errorf("Error while creating the database : %s", err.Error())
	}

	return nil
}
