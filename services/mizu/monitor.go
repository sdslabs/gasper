package mizu

import (
	"math"
	"time"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"

	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

func registerMetrics() {
	apps, err := docker.ListContainers()
	if err != nil {
		utils.LogError(err)
		return
	}
	var parsedMetricsList []interface{}

	for _, app := range apps {
		metrics, err := docker.ContainerStats(app)
		if err != nil {
			utils.LogError(err)
			continue
		}

		containerStatus, err := docker.InspectContainerState(app)
		if err != nil {
			utils.LogError(err)
			continue
		}

		// memory metrics
		memoryUsage := metrics.Memory.Usage
		maxUsage := metrics.Memory.MaxUsage
		memoryLimit := metrics.Memory.Limit
		if memoryLimit == 0 {
			// utils.Log("Container has stopped", utils.ErrorTAG)
			// error needs to be handled in a better way
			continue
		}

		// cpu metrics
		cpuTime := metrics.CPU.CPUUsage.TotalUsage
		onlineCPUs := metrics.CPU.OnlineCPUs
		if onlineCPUs == 0 {
			// utils.Log("Container has stopped", utils.ErrorTAG)
			// error needs to be handled in a better way
			continue
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
		}

		parsedMetricsList = append(parsedMetricsList, parsedMetrics)
	}

	if _, err = mongo.BulkRegisterMetrics(parsedMetricsList); err != nil {
		utils.Log("Failed to register metrics", utils.ErrorTAG)
		utils.LogError(err)
	}
}

// ScheduleMetricsCollection runs the registerMetricsHandler at the given metrics interval
func ScheduleMetricsCollection() {
	interval := configs.ServiceConfig.Mizu.MetricsInterval * time.Second
	scheduler := utils.NewScheduler(interval, registerMetrics)
	scheduler.RunAsync()
}
