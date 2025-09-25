package utils

import (
	"fmt"

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
	percentages, err := cpu.Percent(500*time.Millisecond, false)
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

// GetCPUTimes returns CPU times (user, system, idle, etc.)
func GetCPUTimes() (cpu.TimesStat, error) {
	times, err := cpu.Times(false)
	if err != nil || len(times) == 0 {
		return cpu.TimesStat{}, err
	}
	return times[0], nil
}

// GetSystemIdlePercentage calculates system idle percentage
func GetSystemIdlePercentage() (float64, error) {
	times, err := GetCPUTimes()
	if err != nil {
		return 0, err
	}
	total := times.User + times.System + times.Idle + times.Nice + times.Iowait + times.Irq + times.Softirq + times.Steal
	if total == 0 {
		return 0, fmt.Errorf("total CPU time is zero")
	}
	return (times.Idle / total) * 100, nil
}
