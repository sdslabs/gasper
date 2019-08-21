package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
)

// UpdateOne updates a document in the mongoDB collection
func UpdateOne(collectionName string, filter bson.M, data bson.M) interface{} {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": data}, nil)
	return res
}

// UpdateInstance is an abstraction over UpdateOne which updates an application in mongoDB
func UpdateInstance(filter bson.M, data bson.M) interface{} {
	return UpdateOne(InstanceCollection, filter, data)
}

// UpdateMany updates multiple documents in the mongoDB collection
func UpdateMany(collectionName string, filter bson.M, data bson.M) interface{} {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.UpdateMany(ctx, filter, bson.M{"$set": data}, nil)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

// UpdateInstances is an abstraction over UpdateMany which updates multiple applications in mongoDB
func UpdateInstances(filter bson.M, data bson.M) interface{} {
	return UpdateMany(InstanceCollection, filter, data)
}
