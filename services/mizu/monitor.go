package mizu

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"

	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
	m "go.mongodb.org/mongo-driver/mongo"
)

func monitorHandler() {
	go registerMetrics()
}

func registerMetrics() {
	apps := mongo.FetchAppInfo(types.M{
		mongo.HostIPKey: utils.HostIP,
	},
	)

	var parsedMetricsList []m.WriteModel

	for _, app := range apps {
		metrics, err := docker.ContainerStats(app["name"].(string))
		if err != nil {
			utils.LogError(err)
			return
		}

		var parsedMetrics types.Stats

		parsedMetrics.Name = app["name"].(string)

		// Grab the read time of the metrics
		parsedMetrics.ReadTime = time.Now().Unix()

		// Check if the instance is dead or alive
		var instanceURL string = app["host_ip"].(string) + ":" + strconv.Itoa(int(app["container_port"].(int32)))
		instanceDead := utils.NotAlive(instanceURL)
		parsedMetrics.Alive = !instanceDead

		// Grab and parse the memory metrics
		memoryUsage := metrics["memory_stats"].(map[string]interface{})["usage"].(float64)
		maxUsage := metrics["memory_stats"].(map[string]interface{})["max_usage"].(float64)
		memoryLimit := metrics["memory_stats"].(map[string]interface{})["limit"].(float64)
		parsedMetrics.MemoryUsage = fmt.Sprintf("%f", (memoryUsage / memoryLimit))
		parsedMetrics.MaxMemoryUsage = fmt.Sprintf("%f", (maxUsage / memoryLimit))

		// Grab and parse the cpu metrics
		cpuTime := metrics["cpu_stats"].(map[string]interface{})["cpu_usage"].(map[string]interface{})["total_usage"].(float64)
		onlineCPUs := metrics["cpu_stats"].(map[string]interface{})["online_cpus"].(float64)
		parsedMetrics.OnlineCPUS = fmt.Sprintf("%f", onlineCPUs)
		parsedMetrics.CPUUsage = fmt.Sprintf("%f", (cpuTime / (1000000000 * onlineCPUs)))

		operation := m.NewUpdateOneModel()
		operation.SetFilter(types.M{
			"read_time": time.Now().Unix(),
		}).SetUpdate(parsedMetrics).SetUpsert(true)

		parsedMetricsList = append(parsedMetricsList, operation)
	}

	_, err := mongo.UpsertMetrics(parsedMetricsList)

	if err != nil && err != mongo.ErrNoDocuments {
		utils.LogError(err)
		return
	}
}

// ScheduleCollectMetrics runs the registerMetricsHandler at the given metrics interval
func ScheduleCollectMetrics() {
	interval := configs.ServiceConfig.Kaze.MetricsInterval * time.Second
	scheduler := utils.NewScheduler(interval, monitorHandler)
	scheduler.RunAsync()
}
