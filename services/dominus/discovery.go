package dominus

import (
	"time"

	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

// exposeService exposes a single microservice along with its apps
func exposeService(service string, config map[string]interface{}) {
	if config["deploy"].(bool) {
		apps := mongo.FetchAppInfo(
			map[string]interface{}{
				"language": service,
				"hostIP":   utils.HostIP,
			},
		)
		count := len(apps)
		err := redis.RegisterService(
			service,
			utils.HostIP+config["port"].(string),
			float64(count),
		)
		if err != nil {
			panic(err)
		}

		payload := make(map[string]interface{})

		for _, app := range apps {
			payload[app["name"].(string)] = utils.HostIP + config["port"].(string)
		}
		err = redis.BulkRegisterApps(payload)
		if err != nil {
			panic(err)
		}
	}
}

// exposeServices exposes the microservices running on a host machine for discovery
func exposeServices() {
	for service, config := range configs.ServiceConfig {
		go exposeService(
			service,
			config.(map[string]interface{}),
		)
	}
}

// ScheduleServiceExposure exposes the application services at regular intervals
func ScheduleServiceExposure() {
	interval := time.Duration(configs.SWSConfig["exposureInterval"].(float64)) * time.Second
	scheduler := utils.NewScheduler(interval, exposeServices)
	scheduler.RunAsync()
}
