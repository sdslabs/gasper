package api

import (
	"fmt"

	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// buildAndRun installs application dependencies and starts the application
func buildAndRun(app types.Application) {
	for _, cmd := range app.GetBuildCommands() {
		_, err := docker.ExecProcess(app.GetContainerID(), []string{"sh", "-c", fmt.Sprintf("%s &> /proc/1/fd/1", cmd)})
		if err != nil {
			utils.LogError("API-Build-And-Run-1", err)
		}
	}
	for _, cmd := range app.GetRunCommands() {
		_, err := docker.ExecDetachedProcess(app.GetContainerID(), []string{"sh", "-c", fmt.Sprintf("%s &> /proc/1/fd/1", cmd)})
		if err != nil {
			utils.LogError("API-Build-And-Run-2", err)
		}
	}
}
