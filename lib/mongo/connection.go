package mongo

import (
	"context"
	"os"
	"time"

	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"

	"github.com/sdslabs/gasper/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
var client, err = mongo.Connect(ctx, options.Client().ApplyURI(configs.MongoConfig.URL))
var link = client.Database(projectDatabase)

func setupAdmin() {
	adminInfo := configs.AdminConfig
	pwd, err := utils.HashPassword(adminInfo.Password)
	if err != nil {
		utils.LogError(err)
		os.Exit(1)
	}
	admin := &types.User{
		Username: adminInfo.Username,
		Email:    adminInfo.Email,
		Password: pwd,
		Admin:    true,
	}
	filter := types.M{EmailKey: adminInfo.Email}
	UpsertUser(filter, admin)
	utils.LogInfo("%s (%s) has been given admin privileges", adminInfo.Username, adminInfo.Email)
}

func setup() {
	time.Sleep(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		utils.Log("MongoDB connection was not established", utils.ErrorTAG)
		utils.LogError(err)
		setup()
	} else {
		utils.LogInfo("MongoDB Connection Established")
		setupAdmin()
	}
}

func init() {
	go setup()
}
