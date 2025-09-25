//go:build darwin

package metrics

import (
	"iDevopzAgent/models"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"

	"github.com/shirou/gopsutil/v3/mem"
)

type DarwinCollector struct{}

func GetCollector() Collector {
	return &DarwinCollector{}
}

func (c *DarwinCollector) MetricsCollect(userID, machineID string) (*models.Metrics, error) {
	// CPU metrics
	cpuPercent, _ := cpu.Percent(time.Second, false)
	cpuPerCore, _ := cpu.Percent(time.Second, true)

	// Memory metrics
	memInfo, _ := mem.VirtualMemory()
	swapInfo, _ := mem.SwapMemory()

	// Disk metrics
	diskInfo, _ := disk.Usage("/")
	diskPartitions := getDarwinDiskPartitions()

	// Host info
	hostInfo, _ := host.Info()

	// System idle (100 - CPU usage)
	systemIdle := 100.0
	if len(cpuPercent) > 0 {
		systemIdle = 100.0 - cpuPercent[0]
	}

	// Determine system status
	status := "up"
	if cpuPercent[0] > 90 || memInfo.UsedPercent > 90 || diskInfo.UsedPercent > 95 {
		status = "critical"
	} else if cpuPercent[0] > 80 || memInfo.UsedPercent > 80 || diskInfo.UsedPercent > 85 {
		status = "trouble"
	} else if cpuPercent[0] < 5 && memInfo.UsedPercent < 10 && diskInfo.UsedPercent < 10 {
		status = "down"
	}

	metrics := &models.Metrics{
		UserID:               userID,
		Hostname:             hostInfo.Hostname,
		MachineID:            machineID,
		CPUPercent:           cpuPercent[0],
		CPUPerCore:           cpuPerCore,
		SystemIdle:           systemIdle,
		MemoryUsed:           memInfo.Used,
		MemoryTotal:          memInfo.Total,
		MemoryFree:           memInfo.Free,
		MemoryPercent:        memInfo.UsedPercent,
		SwapMemoryUsed:       swapInfo.Used,
		SwapMemoryTotal:      swapInfo.Total,
		SwapMemoryFree:       swapInfo.Free,
		SwapMemUsagePercent:  swapInfo.UsedPercent,
		DiskUsed:             diskInfo.Used,
		DiskTotal:            diskInfo.Total,
		Uptime:               hostInfo.Uptime,
		MetricGetTime:        time.Now().UTC().Format(time.RFC3339),
		Status:               status,
		Timestamp:            time.Now().Unix(),
		Os:                   hostInfo.OS,
		Interrupts:           0, // Not easily available on macOS
		ContextSwitches:      0, // Not easily available on macOS
		PagesReads:           0, // Not easily available on macOS
		PagesWrites:          0, // Not easily available on macOS
		PageFaults:           0, // Not easily available on macOS
		DiskPartitions:       diskPartitions,
		OverallDiskReadMBps:  0,   // Would need implementation
		OverallDiskWriteMBps: 0,   // Would need implementation
		OverallDiskReadIOPS:  0,   // Would need implementation
		OverallDiskWriteIOPS: 0,   // Would need implementation
		OverallDiskIOPS:      0,   // Would need implementation
		OverallDiskIdle:      100, // Default to idle
		OverallDiskBusy:      0,   // Default to not busy
	}

	return metrics, nil
}

func getDarwinDiskPartitions() []models.DiskPartition {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return []models.DiskPartition{}
	}

	var diskPartitions []models.DiskPartition
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}

		health := "healthy"
		if usage.UsedPercent > 90 {
			health = "critical"
		} else if usage.UsedPercent > 80 {
			health = "warning"
		}

		diskPartitions = append(diskPartitions, models.DiskPartition{
			Device:        p.Device,
			Mountpoint:    p.Mountpoint,
			Fstype:        p.Fstype,
			Total:         usage.Total,
			Used:          usage.Used,
			UsedPercent:   usage.UsedPercent,
			Healthy:       health,
			ReadBytesSec:  0, // Would need iostat implementation
			WriteBytesSec: 0, // Would need iostat implementation
		})
	}
	return diskPartitions
}
