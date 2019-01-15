package dominus

import (
	"fmt"
	"time"

	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

// inspectInstance checks whether a given instance is alive or not and deletes that instance
// if it is dead
func inspectInstance(service, instance string) {
	if utils.NotAlive(instance) {
		err := redis.RemoveServiceInstance(service, instance)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// removeDeadServiceInstances removes all inactive instances in a given service
func removeDeadServiceInstances(service string) {
	instances, err := redis.FetchServiceInstances(service)
	if err != nil {
		fmt.Println(err)
	}
	for _, instance := range instances {
		go inspectInstance(service, instance)
	}
}

// removeDeadInstances removes all inactive instances in every service
func removeDeadInstances() {
	time.Sleep(5 * time.Second)
	var services = []string{"php", "static"}
	for _, service := range services {
		go removeDeadServiceInstances(service)
	}
}

// ScheduleCleanup runs removeDeadInstances on given intervals of time
func ScheduleCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			removeDeadInstances()
		}
	}()
}
