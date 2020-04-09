package configs

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/types"
)

var (
	// configFile is the main configuration file for gasper
	configFile = flag.String("conf", "config.toml", "location of config file")

	// GasperConfig is parsed data for `configFile`
	GasperConfig = getConfiguration()

	// MongoConfig is the configuration for MongoDB
	MongoConfig = GasperConfig.Mongo

	// RedisConfig is the configuration for Redis
	RedisConfig = GasperConfig.Redis

	// ServiceConfig is the configuration for all services
	ServiceConfig = GasperConfig.Services

	// CloudflareConfig is the configuration for cloudflare services used by gasper
	CloudflareConfig = GasperConfig.Cloudflare

	// ImageConfig is the configuration for the images used by gasper
	ImageConfig = GasperConfig.Images

	// AdminConfig is the configuration for default Gasper admin
	AdminConfig = GasperConfig.Admin

	// JWTConfig is the configuration for json web auth token
	JWTConfig = GasperConfig.JWT

	// ServiceMap is the configuration binding the service name to its
	// deployment status and port
	ServiceMap = map[string]*GenericService{
		types.Kaze: {
			Deploy: ServiceConfig.Kaze.Deploy,
			Port:   ServiceConfig.Kaze.Port,
		},
		types.Mizu: {
			Deploy: ServiceConfig.Mizu.Deploy,
			Port:   ServiceConfig.Mizu.Port,
		},
		types.Iwa: {
			Deploy: ServiceConfig.Iwa.Deploy,
			Port:   ServiceConfig.Iwa.Port,
		},
		types.Enrai: {
			Deploy: ServiceConfig.Enrai.Deploy,
			Port:   ServiceConfig.Enrai.Port,
		}, 
		types.MongoDB: {
			Deploy: ServiceConfig.Kaen.MongoDB.PlugIn && ServiceConfig.Kaen.Deploy,
			Port:   ServiceConfig.Kaen.Port,
		},
		types.Hikari: {
			Deploy: ServiceConfig.Hikari.Deploy,
			Port:   ServiceConfig.Hikari.Port,
		},
		types.MySQL: {
			Deploy: ServiceConfig.Kaen.MySQL.PlugIn && ServiceConfig.Kaen.Deploy,
			Port:   ServiceConfig.Kaen.Port,
		},
		types.PostgreSQL: {
			Deploy: ServiceConfig.Kaen.PostgreSQL.PlugIn && ServiceConfig.Kaen.Deploy,
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
