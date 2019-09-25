package redis

import (
	"strings"

	"github.com/go-redis/redis"
)

func keyNotExists(service, url string) bool {
	_, err := client.ZRank(service, url).Result()
	if err != nil {
		return true
	}
	return false
}

// RegisterService puts a service URL in its respective sorted set
func RegisterService(service, url string, score float64) error {
	_, err := client.ZAdd(service, redis.Z{Score: score, Member: url}).Result()
	return err
}

// IncrementServiceLoad increments the number of apps deployed on a service host by 1
func IncrementServiceLoad(service, url string) error {
	_, err := client.ZIncrBy(service, 1, url).Result()
	return err
}

// DecrementServiceLoad decrements the number of apps deployed on a service host by 1
func DecrementServiceLoad(service, url string) error {
	_, err := client.ZIncrBy(service, -1, url).Result()
	return err
}

// GetLeastLoadedInstances returns the URL of the host currently having the least number
// of apps of a particular service deployed
func GetLeastLoadedInstances(service string, count int64) ([]string, error) {
	data, err := client.ZRangeByScore(
		service,
		redis.ZRangeBy{
			Min:    "-inf",
			Max:    "+inf",
			Offset: 0,
			Count:  count,
		}).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []string{"Empty Set"}, nil
	}
	return data, nil
}

// GetLeastLoadedInstance returns a single instance having least number of apps of a particular service deployed
func GetLeastLoadedInstance(service string) (string, error) {
	instance, err := GetLeastLoadedInstances(service, 1)
	if err != nil {
		return ErrEmptySet, err
	}
	return instance[0], nil
}

// FetchServiceInstances returns all instances of a given service
func FetchServiceInstances(service string) ([]string, error) {
	data, err := client.ZRangeByScore(
		service,
		redis.ZRangeBy{
			Min:    "-inf",
			Max:    "+inf",
			Offset: 0,
		}).Result()
	if err != nil {
		return []string{}, err
	}
	if len(data) == 0 {
		return []string{}, nil
	}
	return data, nil
}

// RemoveServiceInstance removes an instance of a particular service
func RemoveServiceInstance(service, member string) error {
	_, err := client.ZRem(service, member).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetSSHPort returns the port of an instance where its ssh service is deployed
func GetSSHPort(url string) (string, error) {
	data, _, err := client.ZScan(SSHKey, 0, url+":*", 1).Result()
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", nil
	}
	return strings.Split(data[0], ":")[1], nil
}
