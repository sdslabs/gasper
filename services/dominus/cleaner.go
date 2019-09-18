package dominus

import (
	"fmt"
	"time"

	"github.com/sdslabs/SWS/lib/commons"
	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

// rescheduleInstance redeploys down instances on least loaded servers
func rescheduleInstance(apps []map[string]interface{}, service string) {
	for _, app := range apps {
		instanceURLs, err := redis.GetLeastLoadedInstances(service, 1)
		if err != nil {
			fmt.Println(err)
		}
		app["hostIP"] = instanceURLs[0]
		commons.Deployer(app, instanceURLs[0], service)
	}
}

// inspectInstance checks whether a given instance is alive or not and deletes that instance
// if it is dead
func inspectInstance(service, instance string) {
	if utils.NotAlive(instance) {
		apps := mongo.FetchAppInfo(
			map[string]interface{}{
				"language": service,
				"hostIP":   utils.HostIP,
			},
		)
		err := redis.RemoveServiceInstance(service, instance)
		if err != nil {
			utils.LogError(err)
		}
		go rescheduleInstance(apps, service)
	}
}

// removeDeadServiceInstances removes all inactive instances in a given service
func removeDeadServiceInstances(service string) {
	instances, err := redis.FetchServiceInstances(service)
	if err != nil {
		utils.LogError(err)
	}
	for _, instance := range instances {
		go inspectInstance(service, instance)
	}
}

// removeDeadInstances removes all inactive instances in every service
func removeDeadInstances() {
	time.Sleep(5 * time.Second)
	for service := range configs.ServiceConfig {
		go removeDeadServiceInstances(service)
	}
}

// ScheduleCleanup runs removeDeadInstances on given intervals of time
func ScheduleCleanup() {
	interval := time.Duration(configs.CronConfig["cleanupInterval"].(float64)) * time.Second
	scheduler := utils.NewScheduler(interval, removeDeadInstances)
	scheduler.RunAsync()
}
