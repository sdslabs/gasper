package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FetchDocs is a generic function which takes a collection name and mongoDB filter as input and returns documents
func FetchDocs(collectionName string, filter types.M, opts ...*options.FindOptions) []types.M {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var data []types.M

	cur, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		utils.LogError(err)
		return nil
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result types.M
		if err := cur.Decode(&result); err != nil {
			utils.LogError(err)
			return nil
		}
		data = append(data, result)
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
		NameKey:         name,
		InstanceTypeKey: AppInstance,
	}).Decode(appCfg)

	return appCfg, err
}

// FetchSingleDatabase returns a database based on a name based filter
func FetchSingleDatabase(name string) (*types.DatabaseConfig, error) {
	collection := link.Collection(InstanceCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	databaseCfg := &types.DatabaseConfig{}

	err := collection.FindOne(ctx, types.M{
		NameKey:         name,
		InstanceTypeKey: DBInstance,
	}).Decode(databaseCfg)

	return databaseCfg, err
}

// FetchSingleUser returns a user based on a email based filter
func FetchSingleUser(email string, opts ...*options.FindOneOptions) (*types.User, error) {
	collection := link.Collection(UserCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &types.User{}
	err := collection.FindOne(ctx, types.M{EmailKey: email}, opts...).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FetchSingleUserWithoutPassword returns a user based on a email based filter without his/her password
func FetchSingleUserWithoutPassword(email string) (*types.User, error) {
	return FetchSingleUser(
		email,
		&options.FindOneOptions{
			Projection: types.M{PasswordKey: 0},
		})
}

// FetchInstanceField returns the value of a given field from an instance
func FetchInstanceField(name, instanceType, field string) (interface{}, error) {
	collection := link.Collection(InstanceCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var data types.M
	err := collection.FindOne(
		ctx,
		types.M{NameKey: name, InstanceTypeKey: instanceType},
		&options.FindOneOptions{
			Projection: types.M{field: 1},
		}).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data[field], nil
}

// FetchDatabaseLanguage returns the language of a database
func FetchDatabaseLanguage(name string) (string, error) {
	data, err := FetchInstanceField(name, DBInstance, LanguageKey)
	if err != nil || data == nil {
		return "", err
	}
	if res, ok := data.(string); ok {
		return res, nil
	}
	return "", errors.New("Failed fetching database language")
}

// FetchDBInfo is an abstraction over FetchDocs for retrieving database related documents
func FetchDBInfo(filter types.M) []types.M {
	filter[InstanceTypeKey] = DBInstance
	return FetchInstances(filter)
}

// FetchUserInfo is an abstraction over FetchDocs for retrieving user details
func FetchUserInfo(filter types.M) []types.M {
	return FetchDocs(
		UserCollection,
		filter,
		&options.FindOptions{
			Projection: types.M{PasswordKey: 0},
		})
}

// FetchContainerMetrics is an abstraction over FetchDocs for retrieving metrics of a container
func FetchContainerMetrics(filter types.M, count int64) []types.M {
	options := options.Find().SetSort(types.M{TimestampKey: -1}).SetLimit(count)
	return FetchDocs(MetricsCollection, filter, options)
}

// CountDocs returns the number of documents matching a filter
func CountDocs(collectionName string, filter types.M) (int64, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return collection.CountDocuments(ctx, filter)
}

// CountInstances returns the number of instances matching a filter
func CountInstances(filter types.M) (int64, error) {
	return CountDocs(InstanceCollection, filter)
}

// CountUsers returns the number of users matching a filter
func CountUsers(filter types.M) (int64, error) {
	return CountDocs(UserCollection, filter)
}

// CountServiceInstances returns the number of applications of a given service deployed
// in a host machine
func CountServiceInstances(service, hostIP string) (int64, error) {
	filter := types.M{
		LanguageKey: service,
		HostIPKey:   hostIP,
	}
	return CountDocs(InstanceCollection, filter)
}
