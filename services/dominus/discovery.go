package dominus

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/sdslabs/SWS/configs"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
)

var instanceRegistrationBindings = map[string]func(instances []map[string]interface{}, currentIP string, config map[string]interface{}){
	"php":     registerApps,
	"node":    registerApps,
	"python":  registerApps,
	"static":  registerApps,
	"mysql":   registerDatabases,
	"mongodb": registerDatabases,
}

func registerApps(instances []map[string]interface{}, currentIP string, config map[string]interface{}) {
	payload := make(map[string]interface{})

	for _, instance := range instances {
		appBind := &types.AppBindings{
			Node:   currentIP + config["port"].(string),
			Server: currentIP + ":" + strconv.FormatInt(int64(instance["httpPort"].(int32)), 10),
		}

		appBindingJSON, err := json.Marshal(appBind)

		if err != nil {
			utils.LogError(err)
			return
		}
		payload[instance["name"].(string)] = appBindingJSON
	}
	err := redis.BulkRegisterApps(payload)
	if err != nil {
		utils.LogError(err)
		return
	}
}

func registerDatabases(instances []map[string]interface{}, currentIP string, config map[string]interface{}) {
	payload := make(map[string]interface{})

	for _, instance := range instances {
		key := fmt.Sprintf("%s:%s", instance["user"].(string), instance["name"].(string))
		payload[key] = currentIP + config["port"].(string)
	}
	err := redis.BulkRegisterDatabases(payload)
	if err != nil {
		utils.LogError(err)
		return
	}
}

// exposeService exposes a single microservice along with its apps
func exposeService(service, currentIP string, config map[string]interface{}) {
	if config["deploy"].(bool) {
		instances := mongo.FetchAppInfo(
			map[string]interface{}{
				"language": service,
				"hostIP":   currentIP,
			},
		)
		count := len(instances)
		err := redis.RegisterService(
			service,
			currentIP+config["port"].(string),
			float64(count),
		)
		if err != nil {
			utils.LogError(err)
			return
		}
		if instanceRegistrationBindings[service] != nil {
			instanceRegistrationBindings[service](instances, currentIP, config)
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
