
package redis

import (
	"fmt"
	"strings"
)

// RegisterDB registers the database in the databases HashMap with its url
func RegisterDB(db, url string) error {
	_, err := client.HSet(DatabaseKey, db, url).Result()
	return err
}

// FetchDBURL returns the URL of the machine where the db in query is deployed
func FetchDBURL(db string) (string, error) {
	result, err := client.HGet(DatabaseKey, db).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// RemoveDB removes the databases's entry from Redis
<<<<<<< HEAD
func RemoveDB(db string) error {
	_, err := client.HDel(DatabaseKey, db).Result()
=======
func RemoveDB(dbKey string) error {
	dbName := strings.Split(dbKey, ":")[0]
	dbUser := strings.Split(dbKey, ":")[1]
	dbKey = fmt.Sprintf(`%s:%s`, dbName, dbUser)
	_, err := client.HDel(DatabaseKey, dbKey).Result()
>>>>>>> checked working
	if err != nil {
		return err
	}
	return nil
}

// FetchAllDatabases gets all the apps with their URL (IP of the node and port)
func FetchAllDatabases() (map[string]string, error) {
	return client.HGetAll(DatabaseKey).Result()
}
