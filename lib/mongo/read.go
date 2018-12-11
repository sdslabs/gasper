package mongo

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
)

// FetchDocs is a generic function which takes a collection name and mongoDB filter as input and returns documents
func FetchDocs(collectionName string, filter bson.M) []map[string]interface{} {
	collection := link.Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var data []map[string]interface{}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		data = append(data, result)
		if err != nil {
			panic(err)
		}
	}
	if err := cur.Err(); err != nil {
		panic(err)
	}
	return data
}

// FetchAppInfo is an abstraction over FetchDocs for retrieving application related documents
func FetchAppInfo(filter bson.M) []map[string]interface{} {
	return FetchDocs("apps", filter)
}
