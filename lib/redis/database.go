package redis

import (
	"encoding/json"

	"github.com/sdslabs/gasper/types"
)

// RegisterDB registers the database in the databases HashMap with its server and node url
func RegisterDB(dbName, nodeURL, serverURL string) error {
	dbBind := &types.InstanceBindings{
		Node:   nodeURL,
		Server: serverURL,
	}
	dbBindingJSON, err := json.Marshal(dbBind)
	if err != nil {
		return err
	}
	_, err = client.HSet(DatabaseKey, dbName, dbBindingJSON).Result()
	return err
}

// FetchDbServer returns the URL of the database's server
func FetchDbServer(dbName string) (string, error) {
	return fetchServer(DatabaseKey, dbName)
}

// FetchDbNode returns the URL of the node where the database is deployed
func FetchDbNode(dbName string) (string, error) {
	return fetchNode(DatabaseKey, dbName)
}

// RemoveDB removes the databases's entry from Redis
func RemoveDB(dbName string) error {
	_, err := client.HDel(DatabaseKey, dbName).Result()
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
