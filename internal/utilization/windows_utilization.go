//go:build windows
// +build windows

package utilization

import (
	"fmt"
	"iDevopzAgent/internal/utils"
	"iDevopzAgent/models"
	"math"
)

type WindowsCollector struct{}

func (l WindowsCollector) CpuUtilization(userID string, machineId string) (*models.CpuUtilization, error) {
	// Get CPU percentage
	cpuUsageRaw, err := utils.GetCPUPercentage()
	if err != nil {
		return nil, err
	}
	// Round to 2 decimal places
	cpuUsage := math.Round(cpuUsageRaw*100) / 100

	hostname, _ := utils.GetHostName()

	return &models.CpuUtilization{
		UserID:     userID,
		MachineID:  machineId,
		Hostname:   hostname,
		CPUPercent: cpuUsage,
	}, nil
}

func (l WindowsCollector) MemoryUtilization(userID string, machineId string) (*models.MemoryUtilization, error) {
	usedPercent, _, _, err := utils.GetMemoryUsage()
	if err != nil {
		return nil, err
	}

	// Round memory usage percentage to 2 decimal places
	memoryUsage := math.Round(usedPercent*100) / 100

	hostname, _ := utils.GetHostName()

	return &models.MemoryUtilization{
		UserID:        userID,
		MachineID:     machineId,
		Hostname:      hostname,
		MemoryPercent: memoryUsage,
	}, nil
}

func (l WindowsCollector) DiskUtilization(userID string, machineId string) (*models.DiskUtilization, error) {
	const path = "/" // Default mount point for Linux

	usedPercent, _, _, _, err := utils.GetDiskUsage(path)
	if err != nil {
		return nil, fmt.Errorf("error getting disk usage: %v", err)
	}

	hostname, _ := utils.GetHostName()

	// Round values
	usedPercent = math.Round(usedPercent*100) / 100
	freePercent := math.Round((100-usedPercent)*100) / 100

	return &models.DiskUtilization{
		UserID:      userID,
		MachineID:   machineId,
		Hostname:    hostname,
		UsedPercent: usedPercent,
		FreePercent: freePercent,
	}, nil
}

func UtilizationCollector() Collector {
	return WindowsCollector{}
}
