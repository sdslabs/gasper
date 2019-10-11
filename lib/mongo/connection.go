package mongo

import (
	"context"
	"time"

	"github.com/sdslabs/SWS/lib/utils"

	"github.com/sdslabs/SWS/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
var client, err = mongo.Connect(ctx, options.Client().ApplyURI(configs.MongoConfig.URL))
var link = client.Database("sws")

func setupAdmin() {
	adminInfo := configs.AdminConfig
	pwd, err := utils.HashPassword(adminInfo["password"].(string))
	if err != nil {
		utils.LogError(err)
		panic(err)
	}
	admin := map[string]interface{}{
		"email":    adminInfo["email"],
		"username": adminInfo["username"],
		"password": pwd,
		"is_admin": true,
	}
	filter := map[string]interface{}{"email": adminInfo["email"]}
	UpsertUser(filter, admin)
	utils.LogInfo("%s (%s) has been given admin privileges", adminInfo["username"], adminInfo["email"])
}

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		utils.LogError(err)
		panic(err)
	} else {
		utils.LogInfo("MongoDB Connection Established")
		setupAdmin()
	}
}
