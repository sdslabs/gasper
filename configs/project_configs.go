package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func parseJSON(path string) *GasperCfg {
	config := &GasperCfg{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("File %s does not exist", path))
	}
	err = json.Unmarshal(file, config)
	if err != nil {
		panic(fmt.Sprintf("Invalid %s file", path))
	}
	return config
}

// configFile is the main configuration file for the API
const configFile = "config.json"

// GasperConfig is parsed data for `configFile`
var GasperConfig = parseJSON(configFile)

// MongoConfig is the configuration for MongoDB
var MongoConfig = GasperConfig.Mongo

// RedisConfig is the configuration for Redis
var RedisConfig = GasperConfig.Redis

// ServiceConfig is the configuration for all services
var ServiceConfig = GasperConfig.Services

// FalconConfig is the configuration for all the falcon client services
var FalconConfig = GasperConfig.Falcon

// CronConfig is the configuration for all the daemons managed by SWS
var CronConfig = GasperConfig.Cron

// CloudflareConfig is the configuration for cloudflare services used by SWS
var CloudflareConfig = GasperConfig.Cloudflare

// ImageConfig is the configuration for the images used by SWS
var ImageConfig = GasperConfig.Images

// ServiceMap is the configuration binding the service name to its
// deployment status and port
var ServiceMap = map[string]*GenericService{
	"dominus": &ServiceConfig.Dominus,
	"mizu":    &ServiceConfig.Mizu,
	"ssh": &GenericService{
		Deploy: ServiceConfig.SSH.Deploy,
		Port:   ServiceConfig.SSH.Port,
	},
	"ssh_proxy": &GenericService{
		Deploy: ServiceConfig.SSHProxy.Deploy,
		Port:   ServiceConfig.SSHProxy.Port,
	},
	"enrai": &ServiceConfig.Enrai,
	"mysql": &GenericService{
		Deploy: ServiceConfig.Mysql.Deploy,
		Port:   ServiceConfig.Mysql.Port,
	},
	"mongodb": &GenericService{
		Deploy: ServiceConfig.Mongodb.Deploy,
		Port:   ServiceConfig.Mongodb.Port,
	},
}
