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
