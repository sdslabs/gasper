package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/utils"
)

var client = redis.NewClient(&redis.Options{
	Addr:     fmt.Sprintf("%s:%d", configs.RedisConfig.Host, configs.RedisConfig.Port),
	Password: configs.RedisConfig.Password,
	DB:       configs.RedisConfig.DB,
})

func init() {
	_, err := client.Ping().Result()
	if err != nil {
		utils.LogError(err)
		panic(err)
	} else {
		utils.LogInfo("Redis Connection Established")
	}
}
