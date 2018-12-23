package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func init() {
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Redis Connection Established")
	}
}
