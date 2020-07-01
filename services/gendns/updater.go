package gendns

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

func handleError(err error) {
	utils.Log("GenDNS-Updater-1", "Failed to update DNS Record Storage", utils.ErrorTAG)
	utils.LogError("GenDNS-Updater-2", err)
}

// filterValidInstances filters the instances and returns
// valid instances i.e which is in the form of IP:Port
func filterValidInstances(reverseProxyInstances []string) []string {
	filteredInstances := make([]string, 0)
	for _, instance := range reverseProxyInstances {
		if strings.Contains(instance, ":") {
			filteredInstances = append(filteredInstances, instance)
		} else {
			utils.LogError("GenDNS-Updater-3", fmt.Errorf("Instance %s is of invalid format", instance))
		}
	}
	return filteredInstances
}

// Updates the DNS record storage periodically
// It assigns the A records in such a way that the load is
// equally distributed among all available GenProxy Reverse Proxy Instances
func updateStorage() {
	reverseProxyInstances, err := redis.FetchServiceInstances(types.GenProxy)
	if err != nil {
		handleError(err)
		return
	}

	reverseProxyInstances = filterValidInstances(reverseProxyInstances)
	if len(reverseProxyInstances) == 0 {
		utils.Log("GenDNS-Updater-4", "No valid GenProxy instances available", utils.ErrorTAG)
		return
	}

	sort.Strings(reverseProxyInstances)
	updateBody := make(map[string]string)
	instanceNum := len(reverseProxyInstances)

	// Create enrties for applications
	appMap, err := redis.FetchAllApps()
	if err != nil {
		handleError(err)
		return
	}
	apps := utils.GetMapKeys(appMap)
	sort.Strings(apps)

	for index, app := range apps {
		fqdn := fmt.Sprintf("%s.app.%s.", app, configs.GasperConfig.Domain)
		address := strings.Split(reverseProxyInstances[index%instanceNum], ":")[0]
		updateBody[fqdn] = address
	}

	// Create enrties for databases
	dbMap, err := redis.FetchAllDatabases()
	if err != nil {
		handleError(err)
		return
	}

	dbInfoStruct := &types.InstanceBindings{}

	for db, data := range dbMap {
		resultByte := []byte(data)
		if err = json.Unmarshal(resultByte, dbInfoStruct); err != nil {
			handleError(err)
			continue
		}
		if strings.Contains(dbInfoStruct.Server, ":") {
			fqdn := fmt.Sprintf("%s.db.%s.", db, configs.GasperConfig.Domain)
			updateBody[fqdn] = strings.Split(dbInfoStruct.Server, ":")[0]
		}
	}

	// Create entry for Master
	masterFQDN := fmt.Sprintf("%s.%s.", types.Master, configs.GasperConfig.Domain)
	rand.Seed(time.Now().Unix())
	address := strings.Split(reverseProxyInstances[rand.Intn(len(reverseProxyInstances))], ":")[0]
	updateBody[masterFQDN] = address

	storage.Replace(updateBody)
}

// ScheduleUpdate runs updateStorage on given intervals of time
func ScheduleUpdate() {
	interval := configs.ServiceConfig.GenDNS.RecordUpdateInterval * time.Second
	scheduler := utils.NewScheduler(interval, updateStorage)
	scheduler.RunAsync()
}
