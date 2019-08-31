package redis

// RegisterApp registers the app in the apps HashMap with its url
func RegisterApp(appName, url string) error {
	_, err := client.HSet(AppKey, appName, url).Result()
	return err
}

// BulkRegisterApps registers multiple apps at once
func BulkRegisterApps(data map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}
	_, err := client.HMSet(AppKey, data).Result()
	return err
}

// FetchAppURL returns the URL of the machine where the app in query is deployed
func FetchAppURL(appName string) (string, error) {
	result, err := client.HGet(AppKey, appName).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// RemoveApp removes the application's entry from Redis
func RemoveApp(appName string) error {
	_, err := client.HDel(AppKey, appName).Result()
	if err != nil {
		return err
	}
	return nil
}
