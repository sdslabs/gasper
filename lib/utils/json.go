package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func parseJSON(path string) map[string]interface{} {
	var config map[string]interface{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("File %s does not exist", path))
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(fmt.Errorf("Invalid %s file", path))
	}
	return config
}

// configFile is the main configuration file for the API
var configFile = "config.json"

// SWSConfig is parsed data for `configFile`
var SWSConfig = parseJSON(configFile)

// MongoConfig is the configuration for MongoDB
var MongoConfig = SWSConfig["mongo"].(map[string]interface{})

// RedisConfig is the configuration for Redis
var RedisConfig = SWSConfig["redis"].(map[string]interface{})

// ServiceConfig is the configuration for all services
var ServiceConfig = SWSConfig["services"].(map[string]interface{})

// FalconConfig is the configuration for all the falcon client services
var FalconConfig = SWSConfig["falcon"].(map[string]interface{})
