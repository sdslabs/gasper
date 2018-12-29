package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sdslabs/SWS/lib/utils"
)

var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
var client, err = mongo.Connect(ctx, utils.MongoConfig["url"].(string))
var link = client.Database("sws")

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("MongoDB Connection Established")
	}
}
