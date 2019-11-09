package mongo

import (
	"context"
	"time"

	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// FetchDocs is a generic function which takes a collection name and mongoDB filter as input and returns documents
func FetchDocs(collectionName string, filter types.M) []types.M {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var data []types.M

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		utils.LogError(err)
		return nil
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result types.M
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

// FetchInstances is an abstraction over FetchDocs for retrieving any instance documents
func FetchInstances(filter types.M) []types.M {
	return FetchDocs(InstanceCollection, filter)
}

// FetchAppInfo is an abstraction over FetchDocs for retrieving application related documents
func FetchAppInfo(filter types.M) []types.M {
	filter[InstanceTypeKey] = AppInstance
	return FetchInstances(filter)
}

// FetchSingleApp returns an application based on a name based filter
func FetchSingleApp(name string) (*types.ApplicationConfig, error) {
	collection := link.Collection(InstanceCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	appCfg := &types.ApplicationConfig{}

	err := collection.FindOne(ctx, types.M{
		"name":          name,
		InstanceTypeKey: AppInstance,
	}).Decode(appCfg)
	if err != nil {
		return nil, err
	}
	return appCfg, nil
}

// FetchSingleDatabase returns a database based on a name based filter
func FetchSingleDatabase(name string) (*types.DatabaseConfig, error) {
	collection := link.Collection(InstanceCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	databaseCfg := &types.DatabaseConfig{}

	err := collection.FindOne(ctx, types.M{
		"name":          name,
		InstanceTypeKey: DBInstance,
	}).Decode(databaseCfg)
	if err != nil {
		return nil, err
	}
	return databaseCfg, nil
}

// FetchDBInfo is an abstraction over FetchDocs for retrieving database related documents
func FetchDBInfo(filter types.M) []types.M {
	filter[InstanceTypeKey] = DBInstance
	return FetchInstances(filter)
}

// FetchUserInfo is an abstraction over FetchDocs for retrieving user details
func FetchUserInfo(filter types.M) []types.M {
	return FetchDocs(UserCollection, filter)
}

// CountDocs returns the number of documents matching a filter
func CountDocs(collectionName string, filter types.M) (int64, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter)
	return count, err
}

// CountInstances returns the number of instances matching a filter
func CountInstances(filter types.M) (int64, error) {
	return CountDocs(InstanceCollection, filter)
}

// CountServiceInstances returns the number of applications of a given service deployed
// in a host machine
func CountServiceInstances(service, hostIP string) (int64, error) {
	filter := types.M{
		"language": service,
		HostIPKey:  hostIP,
	}
	return CountDocs(InstanceCollection, filter)
}
