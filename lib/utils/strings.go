package utils

import (
	"strconv"
	"github.com/sdslabs/gasper/types"
)

// QueryToFilter filters out queries from the URL parameters
func QueryToFilter(queries map[string][]string) types.M {
	filter := make(types.M)

	for key, value := range queries {
		if(key == "git_url"){
			filter["git.repo_url"] = value[0]
			continue
		}
		if(key == "container_port"){
			containerPort, err := strconv.Atoi(value[0])
			if err != nil {
				return nil
			}
			filter["container_port"] = containerPort
			continue
		}
		filter[key] = value[0]
	}

	return filter
}
