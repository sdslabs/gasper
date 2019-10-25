package utils

// Contains check if an string is present in the given string array
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ToStringSlice converts interface{} to []string
func ToStringSlice(v interface{}) []string {
	var strSlice []string
	for _, val := range v.([]interface{}) {
		strSlice = append(strSlice, val.(string))
	}
	return strSlice
}

// GetMapKeys returns the keys present in a map
func GetMapKeys(data map[string]string) []string {
	keys := make([]string, 0)
	for key := range data {
		keys = append(keys, key)
	}
	return keys
}
