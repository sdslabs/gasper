package redis

import (
	"github.com/go-redis/redis"
)

func keyNotExists(service, url string) bool {
	_, err := client.ZRank(service, url).Result()
	if err != nil {
		return true
	}
	return false
}

// RegisterService puts a service URL in its respective sorted set if it doesn't exist
// for service discovery
func RegisterService(service, url string, score float64) error {
	if keyNotExists(service, url) {
		_, err := client.ZAdd(service, redis.Z{Score: score, Member: url}).Result()
		return err
	}
	return nil
}

// IncrementServiceLoad increments the number of apps deployed on a service host by 1
func IncrementServiceLoad(service, url string) error {
	_, err := client.ZIncrBy(service, 1, url).Result()
	return err
}

// GetLeastLoadedService returns the URL of the host currently having the least number
// of apps of a particular service deployed
func GetLeastLoadedService(service string) string {
	data, err := client.ZRangeByScore(
		service,
		redis.ZRangeBy{
			Min:    "-inf",
			Max:    "+inf",
			Offset: 0,
			Count:  1,
		}).Result()
	if err != nil {
		return err.Error()
	}
	if len(data) == 0 {
		return "Empty Set"
	}
	return data[0]
}
