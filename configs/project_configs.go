package configs

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/types"
)

// configFile is the main configuration file for gasper
const configFile = "config.toml"

var (
	// GasperConfig is parsed data for `configFile`
	GasperConfig = parseTOML(configFile)

	// MongoConfig is the configuration for MongoDB
	MongoConfig = GasperConfig.Mongo

	// RedisConfig is the configuration for Redis
	RedisConfig = GasperConfig.Redis

	// ServiceConfig is the configuration for all services
	ServiceConfig = GasperConfig.Services

	// FalconConfig is the configuration for all the falcon client services
	FalconConfig = GasperConfig.Falcon

	// CloudflareConfig is the configuration for cloudflare services used by gasper
	CloudflareConfig = GasperConfig.Cloudflare

	// ImageConfig is the configuration for the images used by gasper
	ImageConfig = GasperConfig.Images

	// AdminConfig is the configuration for default Gasper admin
	AdminConfig = GasperConfig.Admin

	// ServiceMap is the configuration binding the service name to its
	// deployment status and port
	ServiceMap = map[string]*GenericService{
		types.Kaze: &GenericService{
			Deploy: ServiceConfig.Kaze.Deploy,
			Port:   ServiceConfig.Kaze.Port,
		},
		types.Mizu: &ServiceConfig.Mizu,
		types.Iwa: &GenericService{
			Deploy: ServiceConfig.Iwa.Deploy,
			Port:   ServiceConfig.Iwa.Port,
		},
		types.Enrai: &GenericService{
			Deploy: ServiceConfig.Enrai.Deploy,
			Port:   ServiceConfig.Enrai.Port,
		},
		types.Hikari: &GenericService{
			Deploy: ServiceConfig.Hikari.Deploy,
			Port:   ServiceConfig.Hikari.Port,
		},
		types.MySQL: &GenericService{
			Deploy: ServiceConfig.Kaen.MySQL.PlugIn && ServiceConfig.Kaen.Deploy,
			Port:   ServiceConfig.Kaen.Port,
		},
		types.MongoDB: &GenericService{
			Deploy: ServiceConfig.Kaen.MongoDB.PlugIn && ServiceConfig.Kaen.Deploy,
			Port:   ServiceConfig.Kaen.Port,
		},
	}
)

func init() {
	if GasperConfig.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
