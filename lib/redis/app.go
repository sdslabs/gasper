package redis

// RegisterApp registers the app in the Apps HashMap with its url
func RegisterApp(appName, url string) error {
	_, err := client.HSet("Apps", appName, url).Result()
	return err
}

// FetchAppURL returns the URL of the machine where the app in query is deployed
func FetchAppURL(appName string) string {
	result, err := client.HGet("Apps", appName).Result()
	if err != nil {
		return err.Error()
	}
	return result
}
