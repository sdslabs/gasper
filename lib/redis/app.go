package redis

import (
	"encoding/json"

	"github.com/sdslabs/gasper/lib/types"
)

// RegisterApp registers the app in the apps HashMap with its url
func RegisterApp(appName, nodeURL, serverURL string) error {
	appBind := &types.AppBindings{
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

// fetchAppBindings returns a struct containing both the server and node URL
func fetchAppBindings(appName string) (*types.AppBindings, error) {
	result, err := client.HGet(AppKey, appName).Result()
	if err != nil {
		return nil, err
	}

	appInfoStruct := &types.AppBindings{}
	resultByte := []byte(result)
	err = json.Unmarshal(resultByte, appInfoStruct)
	if err != nil {
		return nil, err
	}
	return appInfoStruct, nil
}

// FetchAppServer returns the URL of deployed application bound to the container
func FetchAppServer(appName string) (string, error) {
	app, err := fetchAppBindings(appName)
	if err != nil {
		return "", err
	}
	return app.Server, nil
}

// FetchAppURL returns the URL of the machine where the app in query is deployed
func FetchAppURL(app string) (string, error) {
	result, err := client.HGet(AppKey, app).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// FetchAppNode returns the URL of the machine where the app in query is deployed
func FetchAppNode(appName string) (string, error) {
	app, err := fetchAppBindings(appName)
	if err != nil {
		return "", err
	}
	return app.Node, nil
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
