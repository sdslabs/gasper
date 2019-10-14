package mongo

import (
	"context"
	"time"

	"github.com/sdslabs/gasper/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// DeleteOne deletes a document from a mongoDB collection
func DeleteOne(collectionName string, filter bson.M) (interface{}, error) {
	collection := link.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		utils.LogError(err)
		return nil, err
	}
	return res, nil
}

// DeleteInstance is an abstraction over DeleteOne which deletes an application from mongoDB
func DeleteInstance(filter bson.M) (interface{}, error) {
	return DeleteOne(InstanceCollection, filter)
}

// DeleteUser is an abstraction over DeleteOne which deletes a user from mongoDB
func DeleteUser(filter bson.M) (interface{}, error) {
	return DeleteOne(UserCollection, filter)
}
