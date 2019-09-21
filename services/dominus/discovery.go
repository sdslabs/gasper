package dominus

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
)

// exposeService exposes a single microservice along with its apps
func exposeService(service, currentIP string, config map[string]interface{}) {
	if config["deploy"].(bool) {
		apps := mongo.FetchAppInfo(
			map[string]interface{}{
				"language": service,
				"hostIP":   currentIP,
			},
		)
		count := len(apps)
		err := redis.RegisterService(
			service,
			currentIP+config["port"].(string),
			float64(count),
		)
		if err != nil {
			utils.LogError(err)
			panic(err)
		}

		payload := make(map[string]interface{})

		for _, app := range apps {

			appBind := &types.AppBindings{
				Node:   currentIP + config["port"].(string),
				Server: currentIP + ":" + strconv.FormatInt(app["httpPort"].(int64), 10),
			}

			appBindingJSON, err := json.Marshal(appBind)

			if err != nil {
				utils.LogError(err)
				panic(err)
			}
			payload[app["name"].(string)] = appBindingJSON
		}
		err = redis.BulkRegisterApps(payload)
		if err != nil {
			utils.LogError(err)
			panic(err)
		}
	}
}

// exposeServices exposes the microservices running on a host machine for discovery
func exposeServices() {
	currIP, err := utils.GetOutboundIP()
	if err != nil {
		return
	}
	checkAndUpdateState(currIP)
	for service, config := range configs.ServiceConfig {
		go exposeService(
			service,
			currIP,
			config.(map[string]interface{}),
		)
	}
}

// ScheduleServiceExposure exposes the application services at regular intervals
func ScheduleServiceExposure() {
	interval := time.Duration(configs.CronConfig["exposureInterval"].(float64)) * time.Second
	scheduler := utils.NewScheduler(interval, exposeServices)
	scheduler.RunAsync()
}
