//go:build windows
// +build windows

package metrics

import (
	"iDevopzAgent/internal/utils"
	"iDevopzAgent/models"
	"time"

	"github.com/yusufpapurcu/wmi"
)

type WindowsCollector struct{}

func (w WindowsCollector) MetricsCollect(userID string, machineId string) (*models.Metrics, error) {

	const path = "C:\\"
	osName := utils.GetOS()
	// Get memory usage
	memUsagePercent, memTotal, memUsed, err := utils.GetMemoryUsage()
	if err != nil {
		return nil, err
	}

	//swap memory
	swapMemUsagePercent, swapMemTotal, swapMemUsed, err := utils.GetSwapUsage()

	// Get per-core CPU usage
	perCoreCPU, err := utils.GetPerCoreCPUPercentage()
	if err != nil {
		return nil, err
	}
	// Get system idle percentage
	idlePercent, err := utils.GetSystemIdlePercentage()
	if err != nil {
		return nil, err
	}
	// System-wide context switches
	type Win32_PerfFormattedData_PerfOS_System struct {
		ContextSwitchesPerSec uint64
	}

	// Processor interrupts
	type Win32_PerfFormattedData_PerfOS_Processor struct {
		Name             string
		InterruptsPerSec uint64
	}
	var sys []Win32_PerfFormattedData_PerfOS_System
	err = wmi.Query("SELECT ContextSwitchesPerSec FROM Win32_PerfFormattedData_PerfOS_System", &sys)
	var contextSwitches uint64
	if err == nil && len(sys) > 0 {
		contextSwitches = sys[0].ContextSwitchesPerSec
	}

	var procs []Win32_PerfFormattedData_PerfOS_Processor
	err = wmi.Query("SELECT Name, InterruptsPerSec FROM Win32_PerfFormattedData_PerfOS_Processor WHERE Name='_Total'", &procs)
	var interrupts uint64
	if err == nil && len(procs) > 0 {
		interrupts = procs[0].InterruptsPerSec
	}

	// Memory pages (In/Out/Fault)
	type Win32_PerfFormattedData_PerfOS_Memory struct {
		PageReadsPerSec  uint64
		PageWritesPerSec uint64
		PageFaultsPerSec uint64
	}
	var memPages []Win32_PerfFormattedData_PerfOS_Memory
	err = wmi.Query("SELECT PageReadsPerSec, PageWritesPerSec, PageFaultsPerSec FROM Win32_PerfFormattedData_PerfOS_Memory", &memPages)
	var pageReads, pageWrites, pageFaults uint64
	if err == nil && len(memPages) > 0 {
		pageReads = memPages[0].PageReadsPerSec
		pageWrites = memPages[0].PageWritesPerSec
		pageFaults = memPages[0].PageFaultsPerSec
	}
	// Get disk usage
	diskUsagePercent, diskTotal, diskUsed, _, err := utils.GetDiskUsage(path)
	if err != nil {
		return nil, err
	}
	//overall disk io
	type Win32_PerfFormattedData_PerfDisk_PhysicalDisk struct {
		Name                 string
		DiskReadBytesPerSec  uint64
		DiskWriteBytesPerSec uint64
		DiskReadsPerSec      uint64
		DiskWritesPerSec     uint64
		PercentIdleTime      uint64
	}
	// Overall Disk IO (_Total across all disks)
	var overallDisk []Win32_PerfFormattedData_PerfDisk_PhysicalDisk
	err = wmi.Query("SELECT Name, DiskReadBytesPerSec, DiskWriteBytesPerSec ,DiskReadsPerSec, DiskWritesPerSec, PercentIdleTime FROM Win32_PerfFormattedData_PerfDisk_PhysicalDisk WHERE Name='_Total'", &overallDisk)

	var overallReadMBps, overallWriteMBps float64
	var overallReadIOPS, overallWriteIOPS, overallIOPS uint64
	var overallIdlePercent, overallBusyPercent float64

	if err == nil && len(overallDisk) > 0 {
		overallReadMBps = float64(overallDisk[0].DiskReadBytesPerSec) / (1024 * 1024)
		overallWriteMBps = float64(overallDisk[0].DiskWriteBytesPerSec) / (1024 * 1024)

		overallReadIOPS = overallDisk[0].DiskReadsPerSec
		overallWriteIOPS = overallDisk[0].DiskWritesPerSec
		overallIOPS = overallReadIOPS + overallWriteIOPS

		overallIdlePercent = float64(overallDisk[0].PercentIdleTime)
		overallBusyPercent = 100 - overallIdlePercent
	}

	//invitival disk
	diskPartitions, err := GetDiskPartitionsWithIO()
	if err != nil {
		// handle error but continue
		diskPartitions = []models.DiskPartition{}
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
		UserID:               userID,
		Hostname:             hostname,
		MachineID:            machineId,
		CPUPercent:           cpuPercent,
		MemoryUsed:           uint64(bytesToGB(memUsed)),
		MemoryTotal:          memTotal,
		MemoryFree:           uint64(bytesToGB(memTotal - memUsed)),
		SwapMemoryUsed:       uint64(bytesToGB(swapMemUsed)),
		SwapMemoryTotal:      swapMemTotal,
		SwapMemoryFree:       uint64(bytesToGB(swapMemTotal - swapMemUsed)),
		CPUPerCore:           perCoreCPU,
		SystemIdle:           idlePercent,
		MemoryPercent:        memUsagePercent,
		PagesReads:           uint(pageReads),
		PagesWrites:          uint(pageWrites),
		PageFaults:           uint(pageFaults),
		SwapMemUsagePercent:  swapMemUsagePercent,
		DiskUsed:             diskUsed,
		DiskTotal:            diskTotal,
		Uptime:               uptime,
		MetricGetTime:        utcNow,
		Status:               status,
		Os:                   osName,
		Interrupts:           interrupts,
		ContextSwitches:      contextSwitches,
		DiskPartitions:       diskPartitions,
		OverallDiskReadMBps:  overallReadMBps,
		OverallDiskWriteMBps: overallWriteMBps,
		OverallDiskReadIOPS:  overallReadIOPS,
		OverallDiskWriteIOPS: overallWriteIOPS,
		OverallDiskIOPS:      overallIOPS,
		OverallDiskIdle:      overallIdlePercent,
		OverallDiskBusy:      overallBusyPercent,
	}, nil
}
func bytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}
func GetDiskPartitionsWithIO() ([]models.DiskPartition, error) {
	// Get all partitions
	partitions, err := utils.GetDiskPartitions(true)
	if err != nil {
		return nil, err
	}

	var diskPartitions []models.DiskPartition

	for _, p := range partitions {
		// Get disk usage for partition
		UsedPercent, TotalDisk, DiskUsed, Fstype, err := utils.GetDiskUsage(p.Mountpoint)
		if err != nil {
			continue
		}

		// Determine health status
		health := "healthy"
		if UsedPercent > 90 {
			health = "critical"
		} else if UsedPercent > 80 {
			health = "warning"
		}

		// Query read/write bytes/sec via WMI
		type Win32_PerfFormattedData_PerfDisk_LogicalDisk struct {
			Name                 string
			DiskReadBytesPerSec  uint64
			DiskWriteBytesPerSec uint64
		}
		var diskPerf []Win32_PerfFormattedData_PerfDisk_LogicalDisk

		// Remove trailing backslash from device name for WMI query (e.g., "C:")
		deviceName := p.Device
		if len(deviceName) > 1 && deviceName[len(deviceName)-1] == '\\' {
			deviceName = deviceName[:len(deviceName)-1]
		}

		query := "SELECT Name, DiskReadBytesPerSec, DiskWriteBytesPerSec FROM Win32_PerfFormattedData_PerfDisk_LogicalDisk WHERE Name='" + deviceName + "'"
		_ = wmi.Query(query, &diskPerf)

		var readBytesSec, writeBytesSec uint64
		if len(diskPerf) > 0 {
			readBytesSec = diskPerf[0].DiskReadBytesPerSec
			writeBytesSec = diskPerf[0].DiskWriteBytesPerSec
		}

		// Append to results
		diskPartitions = append(diskPartitions, models.DiskPartition{
			Device:        p.Device,
			Mountpoint:    p.Mountpoint,
			Fstype:        Fstype,
			Total:         TotalDisk,
			Used:          DiskUsed,
			UsedPercent:   UsedPercent,
			Healthy:       health,
			ReadBytesSec:  readBytesSec,
			WriteBytesSec: writeBytesSec,
		})
	}

	return diskPartitions, nil
}

func GetCollector() Collector {
	return WindowsCollector{}
}
