package redis

import "github.com/sdslabs/gasper/types"

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
func RemoveDB(db string) error {
	_, err := client.HDel(DatabaseKey, db).Result()
	if err != nil {
		return err
	}
	return nil
}

// FetchAllDatabases gets all the apps with their URL (IP of the node and port)
func FetchAllDatabases() (map[string]string, error) {
	return client.HGetAll(DatabaseKey).Result()
}

// BulkRegisterDatabases registers multiple databases at once
func BulkRegisterDatabases(data types.M) error {
	if len(data) == 0 {
		return nil
	}
	_, err := client.HMSet(DatabaseKey, data).Result()
	return err
}
