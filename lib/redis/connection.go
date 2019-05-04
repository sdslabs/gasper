package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sdslabs/SWS/lib/configs"
)

var client = redis.NewClient(&redis.Options{
	Addr:     configs.RedisConfig["host"].(string) + configs.RedisConfig["port"].(string),
	Password: configs.RedisConfig["password"].(string),
	DB:       int(configs.RedisConfig["DB"].(float64)),
})

func init() {
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Redis Connection Established")
	}
}
