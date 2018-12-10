package mongo

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
)

func RegisterApp(name string, user_id int, github_url string, language string) interface{} {
	collection := link.Collection("apps")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, bson.M{
		"name":       name,
		"user_id":    user_id,
		"github_url": github_url,
		"language":   language,
	})
	if err != nil {
		panic(err)
	}
	return res.InsertedID
}
