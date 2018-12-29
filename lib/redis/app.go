package redis

// RegisterApp registers the app in the apps HashMap with its url
func RegisterApp(appName, url string) error {
	_, err := client.HSet("apps", appName, url).Result()
	return err
}

// FetchAppURL returns the URL of the machine where the app in query is deployed
func FetchAppURL(appName string) (string, error) {
	result, err := client.HGet("apps", appName).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
