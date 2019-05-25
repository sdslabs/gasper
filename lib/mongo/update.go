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

// UpdateApp is an abstraction over UpdateOne which updates an application in mongoDB
func UpdateApp(filter bson.M, data bson.M) interface{} {
	return UpdateOne("apps", filter, data)
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

// UpdateApps is an abstraction over UpdateMany which updates multiple applications in mongoDB
func UpdateApps(filter bson.M, data bson.M) interface{} {
	return UpdateMany("apps", filter, data)
}

// UpdateHostIP updates the application's host IP address
func UpdateHostIP(oldIP, newIP string) interface{} {
	return UpdateApps(
		map[string]interface{}{
			"hostIP": oldIP,
		},
		map[string]interface{}{
			"hostIP": newIP,
		},
	)
}
