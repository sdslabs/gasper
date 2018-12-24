package utils

// Filters out queries from the URL parameters
func QueryToFilter(queries map[string][]string) map[string]interface{} {
	filter := make(map[string]interface{})

	for key, value := range queries {
		filter[key] = value[0]
	}

	return filter
}
