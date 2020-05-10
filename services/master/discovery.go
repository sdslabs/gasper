package master

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

var instanceRegistrationBindings = map[string]func(instances []types.M, currentIP string, config *configs.GenericService){
	types.AppMaker: registerApps,
	types.MySQL:    registerDatabases,
	types.MongoDB:  registerDatabases,
}

var instanceServiceBindings = map[string]func(currentIP, service string) []types.M{
	types.AppMaker: fetchBoundApps,
	types.MySQL:    fetchBoundDatabases,
	types.MongoDB:  fetchBoundDatabases,
}

func fetchBoundApps(currentIP, service string) []types.M {
	return mongo.FetchAppInfo(
		types.M{
			mongo.HostIPKey: currentIP,
		},
	)
}

func fetchBoundDatabases(currentIP, service string) []types.M {
	return mongo.FetchDBInfo(
		types.M{
			mongo.HostIPKey:   currentIP,
			mongo.LanguageKey: service,
		},
	)
}

func registerApps(instances []types.M, currentIP string, config *configs.GenericService) {
	payload := make(types.M)
	for _, instance := range instances {
		appBind := &types.InstanceBindings{
			Node:   fmt.Sprintf("%s:%d", currentIP, config.Port),
			Server: fmt.Sprintf("%s:%v", currentIP, instance[mongo.ContainerPortKey]),
		}
		appBindingJSON, err := json.Marshal(appBind)
		if err != nil {
			utils.LogError(err)
			continue
		}
		payload[instance[mongo.NameKey].(string)] = appBindingJSON
	}
	if err := redis.BulkRegisterApps(payload); err != nil {
		utils.LogError(err)
	}
}

func registerDatabases(instances []types.M, currentIP string, config *configs.GenericService) {
	payload := make(types.M)
	for _, instance := range instances {
		dbBind := &types.InstanceBindings{
			Node:   fmt.Sprintf("%s:%d", currentIP, config.Port),
			Server: fmt.Sprintf("%s:%v", currentIP, instance[mongo.PortKey]),
		}
		dbBindingJSON, err := json.Marshal(dbBind)
		if err != nil {
			utils.LogError(err)
			continue
		}
		payload[instance[mongo.NameKey].(string)] = dbBindingJSON
	}
	if err := redis.BulkRegisterDatabases(payload); err != nil {
		utils.LogError(err)
	}
}

// exposeService exposes a single microservice along with its apps
func exposeService(service, currentIP string, config *configs.GenericService) {
	count := 0
	var instances []types.M
	if instanceServiceBindings[service] != nil {
		instances = instanceServiceBindings[service](currentIP, service)
		count = len(instances)
	}
	err := redis.RegisterService(
		service,
		fmt.Sprintf("%s:%d", currentIP, config.Port),
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

// exposeServices exposes the microservices running on a host machine for discovery
func exposeServices() {
	currIP, err := utils.GetOutboundIP()
	if err != nil {
		return
	}
	checkAndUpdateState(currIP)
	for service, config := range configs.ServiceMap {
		if config.Deploy {
			go exposeService(service, currIP, config)
		}
	}
}

// ScheduleServiceExposure exposes the application services at regular intervals
func ScheduleServiceExposure() {
	interval := configs.ServiceConfig.ExposureInterval * time.Second
	scheduler := utils.NewScheduler(interval, exposeServices)
	scheduler.RunAsync()
}
