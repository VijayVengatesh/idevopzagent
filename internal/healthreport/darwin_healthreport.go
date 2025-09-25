//go:build darwin

package healthreport

import (
	"iDevopzAgent/models"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

type DarwinHealthReportCollector struct{}

func GetHealthReportCollector() Collector {
	return &DarwinHealthReportCollector{}
}

func (h *DarwinHealthReportCollector) GenerateHealthReport(userID, machineID string) (*models.HealthReport, error) {
	// CPU metrics
	cpuPercent, _ := cpu.Percent(time.Second, false)

	// Memory metrics
	memInfo, _ := mem.VirtualMemory()

	// Disk metrics
	diskInfo, _ := disk.Usage("/")

	// Load average
	loadInfo, _ := load.Avg()

	// Determine health status
	status := "healthy"
	var issues []string

	// Check CPU usage
	if cpuPercent[0] > 90 {
		status = "critical"
		issues = append(issues, "High CPU usage")
	} else if cpuPercent[0] > 70 {
		status = "warning"
		issues = append(issues, "Elevated CPU usage")
	}

	// Check memory usage
	if memInfo.UsedPercent > 95 {
		status = "critical"
		issues = append(issues, "Critical memory usage")
	} else if memInfo.UsedPercent > 80 {
		if status != "critical" {
			status = "warning"
		}
		issues = append(issues, "High memory usage")
	}

	// Check disk usage
	if diskInfo.UsedPercent > 95 {
		status = "critical"
		issues = append(issues, "Critical disk usage")
	} else if diskInfo.UsedPercent > 85 {
		if status != "critical" {
			status = "warning"
		}
		issues = append(issues, "High disk usage")
	}

	// Check load average (for systems with load info)
	cpuCount, _ := cpu.Counts(true)
	if loadInfo.Load1 > float64(cpuCount)*2 {
		status = "critical"
		issues = append(issues, "Very high system load")
	} else if loadInfo.Load1 > float64(cpuCount)*1.5 {
		if status != "critical" {
			status = "warning"
		}
		issues = append(issues, "High system load")
	}

	healthReport := &models.HealthReport{
		UserID:        userID,
		MachineID:     machineID,
		Hostname:      "",
		Availability:  status,
		CPUPercent:    cpuPercent[0],
		MemoryPercent: memInfo.UsedPercent,
		DiskPercent:   diskInfo.UsedPercent,
		Downtimes:     0,
		MetricGetTime: time.Now().Format("2006-01-02 15:04:05"),
		SLA:           "99.9",
	}

	return healthReport, nil
}