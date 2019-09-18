package dominus

import (
	"time"

	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

// inspectInstance checks whether a given instance is alive or not and deletes that instance
// if it is dead
func inspectInstance(service, instance string) {
	if utils.NotAlive(instance) {
		err := redis.RemoveServiceInstance(service, instance)
		if err != nil {
			utils.LogError(err)
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
