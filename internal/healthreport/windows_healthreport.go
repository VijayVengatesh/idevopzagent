//go:build windows
// +build windows

package healthreport

import (
	"fmt"
	"iDevopzAgent/internal/utils"
	"iDevopzAgent/models"
	"math"
	"time"
)

type WindowsCollector struct{}

var (
	totalChecks     int
	resourceHealthy int
	totalUptimeOk   int
)

func (l WindowsCollector) GenerateHealthReport(userId string, machineId string) (*models.HealthReport, error) {
	totalChecks++

	// --- Get system values ---
	hostname, err := utils.GetHostName()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %v", err)
	}

	uptimeSeconds, err := utils.GetUptime()
	if err != nil {
		return nil, fmt.Errorf("failed to get uptime: %v", err)
	}
	uptime := time.Duration(uptimeSeconds) * time.Second

	// --- Get metrics ---
	cpuPercent, err := utils.GetCPUPercentage()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %v", err)
	}

	memPercent, _, _, err := utils.GetMemoryUsage()
	if err != nil {
		return nil, err
	}

	const path = "C:\\"

	diskPercent, _, _, _, err := utils.GetDiskUsage(path)
	if err != nil {
		return nil, fmt.Errorf("error getting disk usage: %v", err)
	}

	// --- Evaluation ---
	const threshold = 95.0
	isResourceHealthy := cpuPercent < threshold && memPercent < threshold && diskPercent < threshold
	isUptimeHealthy := uptime > 24*time.Hour

	if isResourceHealthy {
		resourceHealthy++
	}
	if isUptimeHealthy {
		totalUptimeOk++
	}

	availability := (float64(totalUptimeOk) / float64(totalChecks)) * 100
	sla := (float64(resourceHealthy) / float64(totalChecks)) * 100
	utcNow := time.Now().UTC().Format(time.RFC3339)

	// --- Build Report ---
	report := &models.HealthReport{
		UserID:        userId,
		MachineID:     machineId,
		Hostname:      hostname,
		Availability:  fmt.Sprintf("%.1f %%", availability),
		CPUPercent:    math.Round(cpuPercent*100) / 100,
		MemoryPercent: math.Round(memPercent*100) / 100,
		DiskPercent:   math.Round(diskPercent*100) / 100,
		Downtimes:     totalChecks - resourceHealthy,
		MetricGetTime: utcNow,
		SLA:           fmt.Sprintf("%.2f %%", sla),
	}

	return report, nil
}

func GetHealthReportCollector() Collector {
	return WindowsCollector{}
}
