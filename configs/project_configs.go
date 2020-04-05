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
		types.Master: {
			Deploy: ServiceConfig.Master.Deploy,
			Port:   ServiceConfig.Master.Port,
		},
		types.AppMaker: {
			Deploy: ServiceConfig.AppMaker.Deploy,
			Port:   ServiceConfig.AppMaker.Port,
		},
		types.GenSSH: {
			Deploy: ServiceConfig.GenSSH.Deploy,
			Port:   ServiceConfig.GenSSH.Port,
		},
		types.GenProxy: {
			Deploy: ServiceConfig.GenProxy.Deploy,
			Port:   ServiceConfig.GenProxy.Port,
		},
		types.MongoDB: {
			Deploy: ServiceConfig.DbMaker.MongoDB.PlugIn && ServiceConfig.DbMaker.Deploy,
			Port:   ServiceConfig.DbMaker.Port,
		},
		types.GenDNS: {
			Deploy: ServiceConfig.GenDNS.Deploy,
			Port:   ServiceConfig.GenDNS.Port,
		},
		types.Jikan: {
			Deploy: ServiceConfig.Jikan.Deploy,
			Port:   ServiceConfig.Jikan.Port,
		},
		types.MySQL: {
			Deploy: ServiceConfig.DbMaker.MySQL.PlugIn && ServiceConfig.DbMaker.Deploy,
			Port:   ServiceConfig.DbMaker.Port,
		},
		types.PostgreSQL: {
			Deploy: ServiceConfig.DbMaker.PostgreSQL.PlugIn && ServiceConfig.DbMaker.Deploy,
			Port:   ServiceConfig.DbMaker.Port,
		},
		types.Redis: {
			Deploy: ServiceConfig.DbMaker.Redis.PlugIn && ServiceConfig.DbMaker.Deploy,
			Port:   ServiceConfig.DbMaker.Port,
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
