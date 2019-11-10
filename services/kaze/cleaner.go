package kaze

import (
	"fmt"
	"strings"
	"time"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/commons"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// rescheduleInstance redeploys down instances on least loaded servers
func rescheduleInstance(apps []types.M, service string) {
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
			if !strings.Contains(instance, ":") {
				utils.LogError(fmt.Errorf("Instance %s is in invalid format", instance))
				return
			}
			instanceIP := strings.Split(instance, ":")[0]
			apps := mongo.FetchAppInfo(types.M{
				mongo.HostIPKey: instanceIP,
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
	for service := range configs.ServiceMap {
		// No need to send TCP health-check probe to DNS microservice (Hikari)
		if service == types.Hikari {
			continue
		}
		go removeDeadServiceInstances(service)
	}
}

// ScheduleCleanup runs removeDeadInstances on given intervals of time
func ScheduleCleanup() {
	time.Sleep(10 * time.Second)
	interval := configs.ServiceConfig.Kaze.CleanupInterval * time.Second
	scheduler := utils.NewScheduler(interval, removeDeadInstances)
	scheduler.RunAsync()
}
