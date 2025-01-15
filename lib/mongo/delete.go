package mongo

import (
	"context"
	"time"

	"github.com/sdslabs/gasper/types"
)

// DeleteOne deletes a document from a mongoDB collection
func DeleteOne(collectionName string, filter types.M) (interface{}, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return collection.DeleteOne(ctx, filter)
}

// DeleteInstance is an abstraction over DeleteOne which deletes an application from mongoDB
func DeleteInstance(filter types.M) (interface{}, error) {
	return DeleteOne(InstanceCollection, filter)
}

// DeleteMany deletes multiple documents from a mongoDB collection
func DeleteMany(collectionName string, filter types.M) (interface{}, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return collection.DeleteMany(ctx, filter)
}

// DeleteInstances is an abstraction over DeleteMany which deletes applications from mongoDB
func DeleteInstances(filter types.M) (interface{}, error) {
	return DeleteMany(InstanceCollection, filter)
}

// DeleteUser is an abstraction over DeleteOne which deletes a user from mongoDB
func DeleteUser(filter types.M) (interface{}, error) {
	return DeleteOne(UserCollection, filter)
}

// DeleteMetrics is an abstraction over DeleteOne which deletes a container metrics from mongoDB
func DeleteMetrics(filter types.M) (interface{}, error) {
	return DeleteOne(MetricsCollection, filter)
}
