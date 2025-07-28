//go:build linux
// +build linux

package metrics

import (
	"iDevopzAgent/internal/utils"
	"iDevopzAgent/models"
	"time"
)

type LinuxCollector struct{}

func (l LinuxCollector) MetricsCollect(userID string) (*models.Metrics, error) {

	const path = "/"
	//Get Os Name

	osName := utils.GetOS()
	// Get memory usage
	memUsagePercent, memTotal, memUsed, err := utils.GetMemoryUsage()
	if err != nil {
		return nil, err
	}

	// Get disk usage
	diskUsagePercent, diskTotal, diskUsed, err := utils.GetDiskUsage(path)
	if err != nil {
		return nil, err
	}

	// Wait and get CPU usage
	cpuPercent, err := utils.GetCPUPercentage()
	if err != nil {
		return nil, err
	}

	// Hostname
	hostname, err := utils.GetHostName()
	if err != nil {
		return nil, err
	}

	// Uptime (using gopsutil directly; optionally abstract later)
	uptime, err := utils.GetUptime()
	if err != nil {
		return nil, err
	}

	// Metrics timestamp
	utcNow := time.Now().UTC().Format(time.RFC3339)

	// Determine system status
	status := "up"
	if cpuPercent > 90 || memUsagePercent > 90 || diskUsagePercent > 95 {
		status = "critical"
	} else if cpuPercent > 80 || memUsagePercent > 80 || diskUsagePercent > 85 {
		status = "trouble"
	} else if cpuPercent < 5 && memUsagePercent < 10 && diskUsagePercent < 10 {
		status = "down"
	}

	return &models.Metrics{
		UserID:        userID,
		Hostname:      hostname,
		CPUPercent:    cpuPercent,
		MemoryUsed:    memUsed,
		MemoryTotal:   memTotal,
		DiskUsed:      diskUsed,
		DiskTotal:     diskTotal,
		Uptime:        uptime,
		MetricGetTime: utcNow,
		Status:        status,
		Os:            osName,
	}, nil
}

func GetCollector() Collector {
	return LinuxCollector{}
}
