package dominus

import (
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
)

// updateHostIP updates the application's host IP address
func updateHostIP(oldIP, currentIP string) (interface{}, error) {
	return mongo.UpdateInstances(
		map[string]interface{}{
			"hostIP": oldIP,
		},
		map[string]interface{}{
			"hostIP": currentIP,
		},
	)
}

// updateState updates the IP address of the machine in the application's context
// and re-registers all the microservices and applications deployed
func updateState(currentIP string) {
	utils.LogInfo(
		"IP address of the machine changed from %s to %s\n",
		utils.HostIP,
		currentIP)

	_, err := updateHostIP(utils.HostIP, currentIP)
	if err != nil {
		utils.LogError(err)
		return
	}
	utils.HostIP = currentIP
}

// checkAndUpdateState checks whether the IP address of the machine has changed or not
func checkAndUpdateState(currentIP string) {
	if utils.HostIP != currentIP {
		updateState(currentIP)
	}
}
