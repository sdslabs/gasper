package redis

import (
	"encoding/json"

	"github.com/sdslabs/gasper/types"
)

// RegisterApp registers the app in the applications HashMap with its server and node url
func RegisterApp(appName, nodeURL, serverURL string) error {
	appBind := &types.InstanceBindings{
		Node:   nodeURL,
		Server: serverURL,
	}
	appBindingJSON, err := json.Marshal(appBind)
	if err != nil {
		return err
	}
	_, err = client.HSet(ApplicationKey, appName, appBindingJSON).Result()
	return err
}

// BulkRegisterApps registers multiple apps at once
func BulkRegisterApps(data types.M) error {
	if len(data) == 0 {
		return nil
	}
	_, err := client.HMSet(ApplicationKey, data).Result()
	return err
}

// FetchAppServer returns the URL of deployed application
func FetchAppServer(appName string) (string, error) {
	return fetchServer(ApplicationKey, appName)
}

// FetchAppNode returns the URL of the node where the application is deployed
func FetchAppNode(appName string) (string, error) {
	return fetchNode(ApplicationKey, appName)
}

// RemoveApp removes the application's entry from Redis
func RemoveApp(appName string) error {
	_, err := client.HDel(ApplicationKey, appName).Result()
	if err != nil {
		return err
	}
	return nil
}

// FetchAllApps returns all applications along with their URLs (IP of the node and port)
func FetchAllApps() (map[string]string, error) {
	return client.HGetAll(ApplicationKey).Result()
}
