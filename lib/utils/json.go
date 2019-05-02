package utils

import (
	"encoding/json"
	"io/ioutil"
)

func parseJSON(path string) map[string]interface{} {
	var config map[string]interface{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic("File config.json does not exist")
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		panic("Invalid config.json file")
	}
	return config
}

// SWSConfig is parsed data for `config.json`
var SWSConfig = parseJSON("config.json")

// MongoConfig is the configuration for MongoDB
var MongoConfig = SWSConfig["mongo"].(map[string]interface{})

// RedisConfig is the configuration for Redis
var RedisConfig = SWSConfig["redis"].(map[string]interface{})

// ServiceConfig is the configuration for all services
var ServiceConfig = SWSConfig["services"].(map[string]interface{})

// FalconConfig is the configuration for all the falcon client services
var FalconConfig = SWSConfig["falcon"].(map[string]interface{})
