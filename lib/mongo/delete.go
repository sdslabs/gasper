package mongo

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
)

// DeleteOne deletes a document from a mongoDB collection
func DeleteOne(collectionName string, filter bson.M) interface{} {
	collection := link.Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		panic(err)
	}
	return res
}

// DeleteApp is an abstraction over DeleteOne which deletes an application from mongoDB
func DeleteApp(filter bson.M) interface{} {
	return DeleteOne("apps", filter)
}
