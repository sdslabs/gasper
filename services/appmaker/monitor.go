package appmaker

import (
	"fmt"
	"math"
	"time"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/database"
	"github.com/sdslabs/gasper/lib/docker"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

func registerMetrics() {
	apps, err := docker.ListContainers()
	if err != nil {
		utils.LogError("AppMaker-Monitor-1", err)
		return
	}
	var parsedMetricsList []interface{}

	for _, app := range apps {
		metrics, err := docker.ContainerStats(app)
		if err != nil {
			utils.LogError("AppMaker-Monitor-2", err)
			continue
		}

		containerStatus, err := docker.InspectContainerState(app)
		if err != nil {
			utils.LogError("AppMaker-Monitor-3", err)
			continue
		}

		// memory metrics
		memoryUsage := metrics.Memory.Usage
		maxUsage := metrics.Memory.MaxUsage
		memoryLimit := metrics.Memory.Limit
		if memoryLimit == 0 {
			utils.Log("AppMaker-Monitor-4", fmt.Sprintf("Container %s has stopped", app), utils.ErrorTAG)
			// error needs to be handled in a better way
			continue
		}

		// cpu metrics
		cpuTime := metrics.CPU.CPUUsage.TotalUsage
		onlineCPUs := metrics.CPU.OnlineCPUs
		if onlineCPUs == 0 {
			utils.Log("AppMaker-Monitor-5", fmt.Sprintf("Container %s has stopped", app), utils.ErrorTAG)
			// error needs to be handled in a better way
			continue
		}
		var logs string
		if app == types.MySQL || app == types.PostgreSQL || app == types.MongoDB {
			logs, err = database.LogDB(app)
			if err != nil {
				utils.LogError("AppMaker-Monitor-12", fmt.Errorf("Error in getting logs of %s:,%s", app, err))
			}
		}
		parsedMetrics := types.Metrics{
			Name:           app,
			Alive:          containerStatus.Running,
			ReadTime:       time.Now().Unix(),
			MemoryUsage:    memoryUsage / memoryLimit,
			MaxMemoryUsage: maxUsage / memoryLimit,
			MemoryLimit:    memoryLimit / math.Pow(1024, 3),
			OnlineCPUs:     onlineCPUs,
			CPUUsage:       cpuTime / (math.Pow(10, 9) * onlineCPUs),
			HostIP:         utils.HostIP,
			Logs:           logs,
		}
		if app == types.MySQL || app == types.PostgreSQL || app == types.MongoDB {
			err = mongo.UpdateOneWithUpsert(mongo.MetricsCollection, types.M{"name": app}, parsedMetrics, options.Update().SetUpsert(true))
			if err != nil {
				utils.LogError("AppMaker-Monitor-13", fmt.Errorf("Error in updating metrics of %s:,%s", app, err))
			}
			continue
		}

		parsedMetricsList = append(parsedMetricsList, parsedMetrics)
	}

	if _, err = mongo.BulkRegisterMetrics(parsedMetricsList); err != nil {
		utils.Log("AppMaker-Monitor-6", "Failed to register metrics", utils.ErrorTAG)
		utils.LogError("AppMaker-Monitor-7", err)
	}
}

// ScheduleMetricsCollection runs the registerMetricsHandler at the given metrics interval
func ScheduleMetricsCollection() {
	interval := configs.ServiceConfig.AppMaker.MetricsInterval * time.Second
	scheduler := utils.NewScheduler(interval, registerMetrics)
	scheduler.RunAsync()
}

// checkContainerHealth checks the health of the containers and restarts the unhealthy ones
func checkContainerHealth(){
	apps := fetchAllApplicationNames()
	for _, app := range apps {
		containerStatus, err := docker.InspectContainerHealth(app)
		if err != nil {
			utils.LogError("AppMaker-Monitor-9", err)
			continue
		}
		// If container is unhealthy, log the error and restart the container
		if containerStatus == docker.Container_Unhealthy{
			utils.Log("AppMaker-Monitor-10", fmt.Sprintf("Container %s has stopped", app), utils.ErrorTAG)
			if err := docker.ContainerRestart(app); err != nil {
				utils.LogError("AppMaker-Monitor-11", err)
			}
		}
	}
}

// ScheduleHealthCheck runs the checkContainerHealthHandler at the given health interval
func ScheduleHealthCheck() {
	interval := configs.ServiceConfig.AppMaker.HealthInterval * time.Second
	scheduler := utils.NewScheduler(interval, checkContainerHealth)
	scheduler.RunAsync()
}