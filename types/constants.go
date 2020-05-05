package types

const (
	// Kaze holds the name of `kaze` microservice
	Kaze = "kaze"

	// Mizu holds the name of `mizu` microservice
	Mizu = "mizu"

	// Kaen holds the name of `kaen` microservice
	Kaen = "kaen"

	// MySQL holds the name of `mysql` component under `kaen`
	MySQL = "mysql"

	// MongoDB holds the name of `mongodb` mongo component under `kaen`
	MongoDB = "mongodb"

	// PostgreSQL holds the name of `postgresql` component under 'kaen'
	PostgreSQL = "postgresql"

	// MongoDBGasper holds the name of `mongodb_gasper` mongo component under `kaze`
	MongoDBGasper = "mongodb_gasper"

	//RedisGasper holds the name of `redis_gasper` redis component under `kaze`
	RedisGasper = "redis_gasper"

	// RedisKaen holds the name of `rediskaen` component under 'kaen'
	RedisKaen = "redis"

	// Iwa holds the name of `iwa` microservice
	Iwa = "iwa"

	// Enrai holds the name of `enrai` microservice
	Enrai = "enrai"

	// Hikari holds the name of `hikari` microservice
	Hikari = "hikari"

	// EnraiSSL holds the name of `enrai` microservice with SSL support
	EnraiSSL = "enrai_ssl"

	// DefaultMemory is the default memory allotted to a container
	DefaultMemory = 0.5

	// DefaultCPUs is the default number of CPUs allotted to a container
	DefaultCPUs = 0.25
)
