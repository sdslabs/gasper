package mongo

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
)

// InsertOne inserts a document into a mongoDB collection
func InsertOne(collectionName string, data bson.M) (interface{}, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

// RegisterApp is an abstraction over InsertOne which inserts application info into mongoDB
func RegisterApp(data bson.M) (interface{}, error) {
	return InsertOne("apps", data)
}
