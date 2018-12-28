package dominus

import (
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

// exposeService exposes a single microservice
func exposeService(service string, config map[string]interface{}) {
	if config["deploy"].(bool) {
		count, err := mongo.CountServiceInstances(
			service,
			utils.HostIP,
		)
		if err != nil {
			panic(err)
		}
		err = redis.RegisterService(
			service,
			utils.HostIP+config["port"].(string),
			float64(count),
		)
		if err != nil {
			panic(err)
		}
	}
}

// ExposeServices exposes the microservices running on a host machine for discovery
func ExposeServices() {
	for service, config := range utils.ServiceConfig {
		go exposeService(
			service,
			config.(map[string]interface{}),
		)
	}
}
