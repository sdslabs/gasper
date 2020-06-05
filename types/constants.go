package types

const (
	// Master holds the name of `master` microservice
	Master = "master"

	// AppMaker holds the name of `appmaker` microservice
	AppMaker = "appmaker"

	// DbMaker holds the name of `dbmaker` microservice
	DbMaker = "dbmaker"

	// MySQL holds the name of `mysql` component under `dbmaker`
	MySQL = "mysql"

	// MongoDB holds the name of `mongodb` mongo component under `dbmaker`
	MongoDB = "mongodb"

	// PostgreSQL holds the name of `postgresql` component under 'dbmaker'
	PostgreSQL = "postgresql"

	// MongoDBGasper holds the name of `mongodb_gasper` mongo component under `master`
	MongoDBGasper = "mongodb_gasper"

	//RedisGasper holds the name of `redis_gasper` redis component under `master`
	RedisGasper = "redis_gasper"

	// Redis holds the name of `redis` component under 'dbmaker'
	Redis = "redis"

	// GenSSH holds the name of `genssh` microservice
	GenSSH = "genssh"

	// GenProxy holds the name of `genproxy` microservice
	GenProxy = "genproxy"

	// GenDNS holds the name of `gendns` microservice
	GenDNS = "gendns"

	// GenProxySSL holds the name of `genproxy` microservice with SSL support
	GenProxySSL = "genproxy_ssl"

	// DefaultMemory is the default memory allotted to a container
	DefaultMemory = 0.5

	// DefaultCPUs is the default number of CPUs allotted to a container
	DefaultCPUs = 0.25

	//SeaweedMaster is the master service for Seaweedfs
	SeaweedMaster = "seaweed_master"

	//SeaweedVolume is the volume service for Seaweedfs
	SeaweedVolume = "seaweed_volume"

	//SeaweedFiler is the filer service for Seaweedfs
	SeaweedFiler = "seaweed_filer"

	//SeaweedCronjob is the cronjob service for Seaweedfs
	SeaweedCronjob = "seaweed_cronjob"

	//SeaweedS3 is the Seaweed service that provides support for AmazonS3
	SeaweedS3 = "seaweed_s3"
)
