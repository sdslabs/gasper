package database

import (
	"context"
	"fmt"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoRootUser     = configs.ServiceConfig.Kaen.MongoDB.Env["MONGO_INITDB_ROOT_USERNAME"].(string)
	mongoRootPassword = configs.ServiceConfig.Kaen.MongoDB.Env["MONGO_INITDB_ROOT_PASSWORD"].(string)
	mongoPort         = configs.ServiceConfig.Kaen.MongoDB.ContainerPort
)

func createConnection(ctx context.Context) (*mongo.Client, error) {
	connectionURI := fmt.Sprintf("mongodb://%s:%s@127.0.0.1:%d/admin", mongoRootUser, mongoRootPassword, mongoPort)
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to mongo: %s", err.Error())
	}

	if err = client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("mongo client couldn't connect with background context: %s", err.Error())
	}
	return client, nil
}

func exec(ctx context.Context, conn *mongo.Database, command interface{}) error {
	v := &(mongo.SingleResult{})
	result := conn.RunCommand(ctx, command)
	return (*result).Decode(v)
}

func createUser(ctx context.Context, conn *mongo.Database, db types.Database) error {
	return exec(
		ctx,
		conn,
		bson.D{
			{Key: "createUser", Value: db.GetUser()},
			{Key: "pwd", Value: db.GetPassword()},
			{Key: "roles", Value: bson.A{
				bson.D{
					{Key: "role", Value: "dbOwner"},
					{Key: "db", Value: db.GetName()},
				},
				"readWrite",
			}},
		})
}

func dropUser(ctx context.Context, conn *mongo.Database, user string) error {
	return exec(ctx, conn, bson.D{{Key: "dropUser", Value: user}})
}

func dropDatabase(ctx context.Context, conn *mongo.Database) error {
	return exec(ctx, conn, bson.D{{Key: "dropDatabase", Value: 1}})
}

func refreshMongoDBUser(ctx context.Context, conn *mongo.Database, db types.Database) error {
	err := dropUser(ctx, conn, db.GetUser())
	if err != nil {
		return fmt.Errorf("Error while deleting the user : %s", err.Error())
	}

	if err = createUser(ctx, conn, db); err != nil {
		return fmt.Errorf("Error while creating the database : %s", err.Error())
	}
	return nil
}

// CreateMongoDB creates a database in the mongodb instance with the given database name, user and password
func CreateMongoDB(db types.Database) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := createConnection(ctx)
	if err != nil {
		return err
	}

	databases, err := client.ListDatabaseNames(ctx, bson.M{"name": db.GetName()})
	if err != nil {
		return fmt.Errorf("Error while creating the database : %s", err.Error())
	}

	if utils.Contains(databases, db.GetName()) {
		return fmt.Errorf("Error while creating the database : Database already Exists")
	}

	conn := client.Database(db.GetName())

	if err = createUser(ctx, conn, db); err != nil {
		utils.LogError(err)
		if err = refreshMongoDBUser(ctx, conn, db); err != nil {
			return fmt.Errorf("Error while creating the database : %s", err.Error())
		}
	}

	return nil
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

	if err = dropDatabase(ctx, conn); err != nil {
		return fmt.Errorf("Error while deleting the database : %s", err.Error())
	}

	if err = dropUser(ctx, conn, databaseName); err != nil {
		return fmt.Errorf("Error while deleting the user : %s", err.Error())
	}
	return nil
}
