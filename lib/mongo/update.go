package mongo

import (
	"context"
	"time"

	"github.com/sdslabs/gasper/types"
	m "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateOne updates a document in the mongoDB collection
func UpdateOne(collectionName string, filter types.M, data interface{}, option *options.FindOneAndUpdateOptions) error {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return collection.FindOneAndUpdate(ctx, filter, types.M{"$set": data}, option).Err()
}

// UpdateInstance is an abstraction over UpdateOne which updates an application in mongoDB
func UpdateInstance(filter types.M, data interface{}) error {
	return UpdateOne(InstanceCollection, filter, data, nil)
}

// UpsertInstance is an abstraction over UpdateOne which updates an application in mongoDB
// or inserts it if the corresponding document doesn't exist
func UpsertInstance(filter types.M, data interface{}) error {
	return UpdateOne(InstanceCollection, filter, data, options.FindOneAndUpdate().SetUpsert(true))
}

// UpdateUser is an abstraction over UpdateOne which updates an application in mongoDB
func UpdateUser(filter types.M, data interface{}) error {
	return UpdateOne(UserCollection, filter, data, nil)
}

// UpsertUser is an abstraction over UpdateOne which updates an application in mongoDB
// or inserts it if the corresponding document doesn't exist
func UpsertUser(filter types.M, data interface{}) error {
	return UpdateOne(UserCollection, filter, data, options.FindOneAndUpdate().SetUpsert(true))
}

// BulkUpsert upserts multiple documents using BulkWrite
func BulkUpsert(collectionName string, data []m.WriteModel, options *options.BulkWriteOptions) (interface{}, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return collection.BulkWrite(ctx, data, options)
}

// UpsertMetrics is an abstraction over BulkUpsert which updates multiple metrics documents in mongoDB
// or inserts them if the corresponding document doesn't exist
func UpsertMetrics(data []m.WriteModel) (interface{}, error) {
	return BulkUpsert(MetricsCollection, data, options.BulkWrite().SetOrdered(false))
}

// UpdateMany updates multiple documents in the mongoDB collection
func UpdateMany(collectionName string, filter types.M, data interface{}) (interface{}, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return collection.UpdateMany(ctx, filter, types.M{"$set": data}, nil)
}

// UpdateInstances is an abstraction over UpdateMany which updates multiple applications in mongoDB
func UpdateInstances(filter types.M, data interface{}) (interface{}, error) {
	return UpdateMany(InstanceCollection, filter, data)
}

func UpdateOneWithUpsert(collectionName string, filter types.M, data interface{}, option *options.UpdateOptions) error {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_,err=collection.UpdateOne(ctx, filter, types.M{"$set": data}, option)
	return err
}