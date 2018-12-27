package mongo

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
)

// UpdateOne updates a document in the mongoDB collection
func UpdateOne(collectionName string, filter bson.M, data bson.M) interface{} {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": data}, nil)
	if err != nil {
		panic(err)
	}
	return res
}

// UpdateApp is an abstraction over UpdateOne which deletes an application from mongoDB
func UpdateApp(filter bson.M, data bson.M) interface{} {
	return UpdateOne("apps", filter, data)
}
