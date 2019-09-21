package redis

import (
	"encoding/json"

	"github.com/sdslabs/SWS/lib/types"
)

// RegisterApp registers the app in the apps HashMap with its url
func RegisterApp(appName, nodeURL, serverURL string) error {
	appBind := &types.AppBinding{
		Node:   nodeURL,
		Server: serverURL,
	}
	appBindingJSON, err := json.Marshal(appBind)
	if err != nil {
		return err
	}
	_, regerr := client.HSet(AppKey, appName, appBindingJSON).Result()
	return regerr
}

// BulkRegisterApps registers multiple apps at once
func BulkRegisterApps(data map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}
	_, err := client.HMSet(AppKey, data).Result()
	return err
}

// FetchAppURL returns a struct containing both the server and node URL
func FetchAppURL(appName string) (*types.AppBinding, error) {
	result, err := client.HGet(AppKey, appName).Result()
	if err != nil {
		return nil, err
	}

	var appInfoStruct *types.AppBinding
	resultByte := []byte(result)
	json.Unmarshal(resultByte, appInfoStruct)

	return appInfoStruct, nil
}

// FetchAppServerURL returns the URL of deployed application bound to the container
func FetchAppServerURL(appName string) (string, error) {
	url, err := FetchAppURL(appName)
	if err != nil {
		return "", err
	}
	return url.Server, nil
}

// FetchAppNodeURL returns the URL of the machine where the app in query is deployed
func FetchAppNodeURL(appName string) (string, error) {
	url, err := FetchAppURL(appName)
	if err != nil {
		return "", err
	}
	return url.Node, nil
}

// RemoveApp removes the application's entry from Redis
func RemoveApp(appName string) error {
	_, err := client.HDel(AppKey, appName).Result()
	if err != nil {
		return err
	}
	return nil
}

// FetchAllApps gets all the apps with their URL (IP of the node and port)
func FetchAllApps() (map[string]string, error) {
	return client.HGetAll(AppKey).Result()
}
