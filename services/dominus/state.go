package dominus

import (
	"fmt"
	"time"

	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/utils"
)

// updateState updates the IP address of the machine in the application's context
// and re-registers all the microservices and applications deployed
func updateState() {
	newHostIP := utils.GetOutboundIP()

	fmt.Printf(
		"IP address of the machine changed from %s to %s\n",
		utils.HostIP,
		newHostIP)

	mongo.UpdateHostIP(utils.HostIP, newHostIP)

	utils.HostIP = newHostIP
	ExposeServices()
}

// checkState checks whether the IP address of the machine has changed or not
func checkState() {
	if utils.HostIP != utils.GetOutboundIP() {
		updateState()
	}
}

// ScheduleStateCheckup runs checkState on given intervals of time
func ScheduleStateCheckup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			checkState()
		}
	}()
}
