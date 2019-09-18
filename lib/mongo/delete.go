package mongo

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/sdslabs/SWS/lib/utils"
)

// DeleteOne deletes a document from a mongoDB collection
func DeleteOne(collectionName string, filter bson.M) interface{} {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		utils.LogError(err)
		panic(err)
	}
	return res
}

// DeleteInstance is an abstraction over DeleteOne which deletes an application from mongoDB
func DeleteInstance(filter bson.M) interface{} {
	return DeleteOne(InstanceCollection, filter)
}
