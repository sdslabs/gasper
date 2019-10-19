package utils

import "github.com/sdslabs/gasper/types"

// QueryToFilter filters out queries from the URL parameters
func QueryToFilter(queries map[string][]string) types.M {
	filter := make(types.M)

	for key, value := range queries {
		filter[key] = value[0]
	}

	return filter
}
