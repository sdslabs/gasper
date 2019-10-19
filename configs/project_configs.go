package configs

const (
	// configFile is the main configuration file for gasper
	configFile = "config.json"

	// Dominus holds the name of `dominus` microservice
	Dominus = "dominus"

	// Mizu holds the name of `mizu` microservice
	Mizu = "mizu"

	// SSH holds the name of `ssh` microservice
	SSH = "ssh"

	// SSHProxy holds the name of `ssh_proxy` microservice
	SSHProxy = "ssh_proxy"

	// Enrai holds the name of `enrai` microservice
	Enrai = "enrai"

	// MySQL holds the name of `mysql` microservice
	MySQL = "mysql"

	// MongoDB holds the name of `mongodb` microservice
	MongoDB = "mongodb"
)

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
		Dominus: &ServiceConfig.Dominus,
		Mizu:    &ServiceConfig.Mizu,
		Ssh: &GenericService{
			Deploy: ServiceConfig.SSH.Deploy,
			Port:   ServiceConfig.SSH.Port,
		},
		SshProxy: &GenericService{
			Deploy: ServiceConfig.SSHProxy.Deploy,
			Port:   ServiceConfig.SSHProxy.Port,
		},
		Enrai: &ServiceConfig.Enrai,
		MySQL: &GenericService{
			Deploy: ServiceConfig.Mysql.Deploy,
			Port:   ServiceConfig.Mysql.Port,
		},
		MongoDB: &GenericService{
			Deploy: ServiceConfig.Mongodb.Deploy,
			Port:   ServiceConfig.Mongodb.Port,
		},
	}
)
