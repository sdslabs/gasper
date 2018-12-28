package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

var client = redis.NewClient(&redis.Options{
	Addr:     utils.RedisConfig["host"].(string) + utils.RedisConfig["port"].(string),
	Password: utils.RedisConfig["password"].(string),
	DB:       utils.RedisConfig["DB"].(int),
})

func init() {
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Redis Connection Established")
	}
}
