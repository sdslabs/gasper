package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/utils"
)

var client = redis.NewClient(&redis.Options{
	Addr:     fmt.Sprintf("%s:%d", configs.RedisConfig.Host, configs.RedisConfig.Port),
	Password: configs.RedisConfig.Password,
	DB:       configs.RedisConfig.DB,
})

func setup() {
	_, err := client.Ping().Result()
	if err != nil {
		utils.Log("Redis-Connection-1", "Redis connection was not established", utils.ErrorTAG)
		utils.LogError("Redis-Connection-2", err)
		time.Sleep(5 * time.Second)
		setup()
	} else {
		utils.LogInfo("Redis-Connection-3", "Redis Connection Established")
	}
}

func init() {
	go setup()
}
