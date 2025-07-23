package utils

import (
	"errors"
	"math"
	"runtime"

	"github.com/shirou/gopsutil/v4/mem"
)

// validateRAMvalue checks if the CPU value entered by the use is within bounds of 0 < value < TotalCPUs-1
func ValidateCPUvalue(cpu float64) (bool, error) {
	totalCPU := runtime.NumCPU()
	if cpu > 0 && cpu <= float64(totalCPU-1) { // 1 CPU reserved for system
		return true, nil
	} else if cpu > float64(totalCPU-1) && cpu <= float64(totalCPU) {
		return false, errors.New("cpu value too high, risks system failure")
	} else {
		return false, errors.New("invalid cpu value")
	}
}

// validateRAMvalue checks if the Memory value entered by the use is within bounds of 0 < value < TotalRAM-1
func ValidateRAMvalue(memory float64) (bool, error) {
	vMemory, err := mem.VirtualMemory()
	if err != nil {
		LogError("Error fetching memory", err)
		return false, errors.New("Unable to assign Memory value")
	}
	totalMemory := float64(vMemory.Total) / math.Pow(1024, 3)
	if memory > 0 && memory <= totalMemory-1 { // 1 GB reserved for system
		return true, nil
	} else if memory > totalMemory-1 && memory <= totalMemory {
		return false, errors.New("memory value too high, risks system failure")
	} else {
		return false, errors.New("invalid memory value")
	}
}
