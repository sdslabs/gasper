package redis

import (
	"fmt"
	"strings"
)

// RegisterDB registers the database in the databases HashMap with its url
func RegisterDB(dbKey, url string) error {
	_, err := client.HSet(DatabaseKey, dbKey, url).Result()
	return err
}

// FetchDBURL returns the URL of the machine where the db in query is deployed
func FetchDBURL(dbKey string) (string, error) {
	result, err := client.HGet(DatabaseKey, dbKey).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// RemoveDB removes the databases's entry from Redis
func RemoveDB(dbKey string) error {
	dbName := strings.Split(dbKey, ":")[0]
	dbUser := strings.Split(dbKey, ":")[1]
	dbKey = fmt.Sprintf(`%s:%s`, dbName, dbUser)
	_, err := client.HDel(DatabaseKey, dbKey).Result()
	if err != nil {
		return err
	}
	return nil
}
