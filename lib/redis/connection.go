package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

var client = redis.NewClient(&redis.Options{
	Addr:     utils.SWSConfig.Redis.Host + utils.SWSConfig.Redis.Port,
	Password: utils.SWSConfig.Redis.Password, // no password set
	DB:       utils.SWSConfig.Redis.DB,       // use default DB
})

func init() {
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Redis Connection Established")
	}
}
