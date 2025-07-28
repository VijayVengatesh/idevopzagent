package utils

import (
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

// GetNumCPU returns the number of logical CPUs
func GetNumCPU() int {
	return runtime.NumCPU()
}

// GetCPUPercentage returns the total CPU usage as a percentage
func GetCPUPercentage() (float64, error) {
	percentages, err := cpu.Percent(0, false)
	if err != nil || len(percentages) == 0 {
		return 0, err
	}
	return percentages[0], nil
}

// GetPerCoreCPUPercentage returns per-core CPU usage percentage
func GetPerCoreCPUPercentage() ([]float64, error) {
	return cpu.Percent(time.Second, true)
}

// GetCPUInfo returns detailed CPU model info
func GetCPUInfo() ([]cpu.InfoStat, error) {
	return cpu.Info()
}
