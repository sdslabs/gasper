package dbmaker

import (
	"fmt"
	"math"
	"time"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"

	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

func registerMetrics() {
	dbs, err := docker.ListContainers()
	if err != nil {
		utils.LogError("DbMaker-Monitor-1", err)
		return
	}
	var parsedMetricsList []interface{}

	for _, db := range dbs {
		metrics, err := docker.ContainerStats(db)
		if err != nil {
			utils.LogError("DbMaker-Monitor-2", err)
			continue
		}

		containerStatus, err := docker.InspectContainerState(db)
		if err != nil {
			utils.LogError("DbMaker-Monitor-3", err)
			continue
		}

		// memory metrics
		memoryUsage := metrics.Memory.Usage
		maxUsage := metrics.Memory.MaxUsage
		memoryLimit := metrics.Memory.Limit
		if memoryLimit == 0 {
			utils.Log("DbMaker-Monitor-4", fmt.Sprintf("Container %s has stopped", db), utils.ErrorTAG)
			// error needs to be handled in a better way
			continue
		}

		// cpu metrics
		cpuTime := metrics.CPU.CPUUsage.TotalUsage
		onlineCPUs := metrics.CPU.OnlineCPUs
		if onlineCPUs == 0 {
			utils.Log("DbMaker-Monitor-5", fmt.Sprintf("Container %s has stopped", db), utils.ErrorTAG)
			// error needs to be handled in a better way
			continue
		}

		parsedMetrics := types.Metrics{
			Name:           db,
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
		utils.Log("DbMaker-Monitor-6", "Failed to register metrics", utils.ErrorTAG)
		utils.LogError("DbMaker-Monitor-7", err)
	}
}

// ScheduleMetricsCollection runs the registerMetricsHandler at the given metrics interval
func ScheduleMetricsCollection() {
	interval := configs.ServiceConfig.DbMaker.MetricsInterval * time.Second
	scheduler := utils.NewScheduler(interval, registerMetrics)
	scheduler.RunAsync()
}
