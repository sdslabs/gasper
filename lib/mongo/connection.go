package mongo

import (
	"context"
	"github.com/sdslabs/SWS/lib/utils"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sdslabs/SWS/lib/configs"
)

var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
var client, err = mongo.Connect(ctx, configs.MongoConfig["url"].(string))
var link = client.Database("sws")

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	} else {
		utils.Log("MongoDB Connection Established")
	}
}
