package hikari

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

func handleError(err error) {
	utils.Log("Failed to update DNS Record Storage", utils.ErrorTAG)
	utils.LogError(err)
}

// filterValidInstances filters the reverse proxy instances and returns
// valid instances i.e which is in the form of IP:Port
func filterValidInstances(reverseProxyInstances []string) []string {
	filteredInstances := make([]string, 0)
	for _, instance := range reverseProxyInstances {
		if strings.Contains(instance, ":") {
			filteredInstances = append(filteredInstances, instance)
		} else {
			utils.LogError(fmt.Errorf("Instance %s in Enrai instances is of invalid format", instance))
		}
	}
	return filteredInstances
}

// Updates the DNS record storage periodically
// It assigns the A records in such a way that the load is
// equally distributed among all available Enrai Reverse Proxy Instances
func updateStorage() {
	appMap, err := redis.FetchAllApps()
	if err != nil {
		handleError(err)
		return
	}
	apps := utils.GetMapKeys(appMap)
	sort.Strings(apps)

	reverseProxyInstances, err := redis.FetchServiceInstances(types.Enrai)
	if err != nil {
		handleError(err)
		return
	}

	reverseProxyInstances = filterValidInstances(reverseProxyInstances)
	if len(reverseProxyInstances) == 0 {
		utils.Log("No valid Enrai instances available", utils.ErrorTAG)
		return
	}

	sort.Strings(reverseProxyInstances)
	updateBody := make(map[string]string)
	instanceNum := len(reverseProxyInstances)

	for index, app := range apps {
		fqdn := fmt.Sprintf("%s.app.%s.", app, configs.GasperConfig.Domain)
		address := strings.Split(reverseProxyInstances[index%instanceNum], ":")[0]
		updateBody[fqdn] = address
	}
	storage.SetBulk(updateBody)
}

// ScheduleUpdate runs updateStorage on given intervals of time
func ScheduleUpdate() {
	interval := configs.ServiceConfig.Hikari.RecordUpdateInterval * time.Second
	scheduler := utils.NewScheduler(interval, updateStorage)
	scheduler.RunAsync()
}
