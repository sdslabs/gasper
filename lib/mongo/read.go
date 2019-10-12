package mongo

import (
	"context"
	"time"

	"github.com/sdslabs/gasper/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// FetchDocs is a generic function which takes a collection name and mongoDB filter as input and returns documents
func FetchDocs(collectionName string, filter bson.M) []map[string]interface{} {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var data []map[string]interface{}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		utils.LogError(err)
		return nil
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		data = append(data, result)
		if err != nil {
			utils.LogError(err)
			return nil
		}
	}
	if err := cur.Err(); err != nil {
		utils.LogError(err)
		return nil
	}
	return data
}

// FetchAppInfo is an abstraction over FetchDocs for retrieving application related documents
func FetchAppInfo(filter bson.M) []map[string]interface{} {
	return FetchDocs(InstanceCollection, filter)
}

// FetchDBInfo is an abstraction over FetchDocs for retrieving database related documents
func FetchDBInfo(filter bson.M) []map[string]interface{} {
	return FetchDocs(InstanceCollection, filter)
}

// FetchDBs is an abstraction over FetchDocs for retrieving details of all the databases
func FetchDBs(filter bson.M) []map[string]interface{} {
	return FetchDocs(InstanceCollection, filter)
}

// FetchUserInfo is an abstraction over FetchDocs for retrieving user details
func FetchUserInfo(filter bson.M) []map[string]interface{} {
	return FetchDocs(UserCollection, filter)
}

// FetchUsers is an abstraction over FetchDocs for retreiving users
func FetchUsers(filter bson.M) []map[string]interface{} {
	return FetchDocs(UserCollection, filter)
}

// CountDocs returns the number of documents matching a filter
func CountDocs(collectionName string, filter bson.M) (int64, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter)
	return count, err
}

// CountInstances returns the number of instances matching a filter
func CountInstances(filter bson.M) (int64, error) {
	return CountDocs(InstanceCollection, filter)
}

// CountServiceInstances returns the number of applications of a given service deployed
// in a host machine
func CountServiceInstances(service, hostIP string) (int64, error) {
	filter := bson.M{
		"language": service,
		"hostIP":   hostIP,
	}
	return CountDocs(InstanceCollection, filter)
}
