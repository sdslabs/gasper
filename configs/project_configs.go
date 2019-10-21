package configs

import "github.com/sdslabs/gasper/types"

// configFile is the main configuration file for gasper
const configFile = "config.json"

var (
	// GasperConfig is parsed data for `configFile`
	GasperConfig = parseJSON(configFile)

	// MongoConfig is the configuration for MongoDB
	MongoConfig = GasperConfig.Mongo

	// RedisConfig is the configuration for Redis
	RedisConfig = GasperConfig.Redis

	// ServiceConfig is the configuration for all services
	ServiceConfig = GasperConfig.Services

	// FalconConfig is the configuration for all the falcon client services
	FalconConfig = GasperConfig.Falcon

	// CronConfig is the configuration for all the daemons managed by gasper
	CronConfig = GasperConfig.Cron

	// CloudflareConfig is the configuration for cloudflare services used by gasper
	CloudflareConfig = GasperConfig.Cloudflare

	// ImageConfig is the configuration for the images used by gasper
	ImageConfig = GasperConfig.Images

	// AdminConfig is the configuration for default Gasper admin
	AdminConfig = GasperConfig.Admin

	// ServiceMap is the configuration binding the service name to its
	// deployment status and port
	ServiceMap = map[string]*GenericService{
		types.Dominus: &ServiceConfig.Dominus,
		types.Mizu:    &ServiceConfig.Mizu,
		types.SSH: &GenericService{
			Deploy: ServiceConfig.SSH.Deploy,
			Port:   ServiceConfig.SSH.Port,
		},
		types.SSHProxy: &GenericService{
			Deploy: ServiceConfig.SSHProxy.Deploy,
			Port:   ServiceConfig.SSHProxy.Port,
		},
		types.Enrai: &GenericService{
			Deploy: ServiceConfig.Enrai.Deploy,
			Port:   ServiceConfig.Enrai.Port,
		},
		types.Hikari: &ServiceConfig.Hikari,
		types.MySQL: &GenericService{
			Deploy: ServiceConfig.Mysql.Deploy,
			Port:   ServiceConfig.Mysql.Port,
		},
		types.MongoDB: &GenericService{
			Deploy: ServiceConfig.Mongodb.Deploy,
			Port:   ServiceConfig.Mongodb.Port,
		},
	}
)
