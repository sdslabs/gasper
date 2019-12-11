package mongo

import (
	"context"
	"time"
)

// InsertOne inserts a document into a mongoDB collection
func InsertOne(collectionName string, data interface{}) (interface{}, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

// RegisterInstance is an abstraction over InsertOne which inserts application info into mongoDB
func RegisterInstance(data interface{}) (interface{}, error) {
	return InsertOne(InstanceCollection, data)
}

// RegisterUser is an abstraction over InsertOne which inserts user into the mongoDB
func RegisterUser(data interface{}) (interface{}, error) {
	return InsertOne(UserCollection, data)
}

// RegisterMetrics is an abstraction over InsertOne which inserts metrics into the mongoDB
func RegisterMetrics(data interface{}) (interface{}, error) {
	return InsertOne(MetricsCollection, data)
}
