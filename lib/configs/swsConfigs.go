package configs

import (
	"encoding/json"
	"io/ioutil"
)

func parseJSON(path string) (map[string]interface{}, error) {
	var config map[string]interface{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// SWSConfig is parsed data for `config.json`
var SWSConfig, _ = parseJSON("config.json")

// MongoConfig is the configuration for MongoDB
var MongoConfig = SWSConfig["mongo"].(map[string]interface{})

// RedisConfig is the configuration for Redis
var RedisConfig = SWSConfig["redis"].(map[string]interface{})

// ServiceConfig is the configuration for all services
var ServiceConfig = SWSConfig["services"].(map[string]interface{})

// FalconConfig is the configuration for all the falcon client services
var FalconConfig = SWSConfig["falcon"].(map[string]interface{})
