package enrai

import (
	"encoding/json"
	"time"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

func handleError(err error) {
	utils.Log("Failed to update Record Storage", utils.ErrorTAG)
	utils.LogError(err)
}

// Updates the reverse proxy record storage periodically
func updateStorage() {
	apps, err := redis.FetchAllApps()
	if err != nil {
		handleError(err)
		return
	}

	updateBody := make(map[string]string)
	appInfoStruct := &types.AppBindings{}

	for name, data := range apps {
		resultByte := []byte(data)
		err = json.Unmarshal(resultByte, appInfoStruct)
		if err != nil {
			handleError(err)
			continue
		}
		updateBody[name] = appInfoStruct.Server
	}
	storage.Replace(updateBody)
}

// ScheduleUpdate runs updateStorage on given intervals of time
func ScheduleUpdate() {
	time.Sleep(10 * time.Second)
	interval := configs.ServiceConfig.Enrai.RecordUpdateInterval * time.Second
	scheduler := utils.NewScheduler(interval, updateStorage)
	scheduler.RunAsync()
}
