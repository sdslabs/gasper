package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func parseJSON(path string) map[string]interface{} {
	var config map[string]interface{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("File %s does not exist", path))
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(fmt.Sprintf("Invalid %s file", path))
	}
	return config
}

// configFile is the main configuration file for the API
const configFile = "config.json"

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

// CronConfig is the configuration for all the daemons managed by SWS
var CronConfig = SWSConfig["cron"].(map[string]interface{})

// CloudflareConfig is the configuration for cloudflare services used by SWS
var CloudflareConfig = SWSConfig["cloudflare"].(map[string]interface{})

// ImageConfig is the configuration for the images used by SWS
var ImageConfig = SWSConfig["images"].(map[string]interface{})

// AdminConfig is the configuration for default SWS admin
var AdminConfig = SWSConfig["admin"].(map[string]interface{})
