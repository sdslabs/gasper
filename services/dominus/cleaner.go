package dominus

import (
	"strings"
	"time"

	"github.com/sdslabs/SWS/configs"
	"github.com/sdslabs/SWS/lib/commons"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

// rescheduleInstance redeploys down instances on least loaded servers
func rescheduleInstance(apps []map[string]interface{}, service string) {
	if len(apps) == 0 {
		return
	}
	for _, app := range apps {
		instanceURL, err := redis.GetLeastLoadedWorker()
		if err != nil {
			utils.LogError(err)
		}
		if instanceURL != redis.ErrEmptySet {
			commons.DeployRPC(app, instanceURL, service)
		}
	}
}

// inspectInstance checks whether a given instance is alive or not and deletes that instance
// if it is dead
func inspectInstance(service, instance string) {
	if utils.NotAlive(instance) {
		err := redis.RemoveServiceInstance(service, instance)
		if err != nil {
			utils.LogError(err)
		}
		if service == "mizu" {
			instanceIP := strings.Split(instance, ":")
			apps := mongo.FetchAppInfo(
				map[string]interface{}{
					"hostIP":       instanceIP[0],
					"instanceType": mongo.AppInstance,
				},
			)
			go rescheduleInstance(apps, service)
		}
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
