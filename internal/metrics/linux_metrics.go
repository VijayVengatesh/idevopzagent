//go:build linux
// +build linux

package metrics

import (
	"bufio"
	"iDevopzAgent/internal/utils"
	"iDevopzAgent/models"
	"os"
	"strconv"
	"strings"
	"time"
)

type LinuxCollector struct{}

func (l LinuxCollector) MetricsCollect(userID string, machineId string) (*models.Metrics, error) {

	const path = "/"
	//Get Os Name

	osName := utils.GetOS()
	// Get memory usage
	memUsagePercent, memTotal, memUsed, err := utils.GetMemoryUsage()
	if err != nil {
		return nil, err
	}

	// Swap memory
	swapMemUsagePercent, swapMemTotal, swapMemUsed, err := utils.GetSwapUsage()
	if err != nil {
		return nil, err
	}

	// Memory pages in/out/fault
	pageReads, pageWrites, pageFaults, err := getLinuxMemoryPages()
	if err != nil {
		return nil, err
	}

	// Get disk usage
	diskUsagePercent, diskTotal, diskUsed, _, err := utils.GetDiskUsage(path)
	if err != nil {
		return nil, err
	}
	// Get overall disk stats
	overallReadMBps, overallWriteMBps, readIOPS, writeIOPS, totalIOPS, diskBusy, diskIdle, err := getOverallDiskStats()
	if err != nil {
		overallReadMBps, overallWriteMBps, readIOPS, writeIOPS, totalIOPS, diskBusy, diskIdle = 0, 0, 0, 0, 0, 0, 100
	}

	// Disk partitions with IO
	diskPartitions, err := GetLinuxDiskPartitionsWithIO()
	if err != nil {
		diskPartitions = []models.DiskPartition{}
	}
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

	// Get interrupts and context switches (from /proc/stat)
	interrupts, contextSwitches, err := getLinuxStats()
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
		OverallDiskReadIOPS:  readIOPS,
		OverallDiskWriteIOPS: writeIOPS,
		OverallDiskIOPS:      totalIOPS,
		OverallDiskIdle:      diskIdle,
		OverallDiskBusy:      diskBusy,
	}, nil
}

func GetCollector() Collector {
	return LinuxCollector{}
}
func bytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}

// Swap memory from /proc/meminfo
func getLinuxSwap() (usedPercent float64, total, used uint64, err error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, 0, err
	}
	defer file.Close()

	var swapTotal, swapFree uint64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		val, _ := strconv.ParseUint(fields[1], 10, 64)
		switch fields[0] {
		case "SwapTotal:":
			swapTotal = val * 1024 // KB to Bytes
		case "SwapFree:":
			swapFree = val * 1024
		}
	}

	if swapTotal > 0 {
		used = swapTotal - swapFree
		usedPercent = (float64(used) / float64(swapTotal)) * 100
	}
	return usedPercent, swapTotal, used, nil
}

// Memory pages in/out/fault from /proc/vmstat
func getLinuxMemoryPages() (pageIn, pageOut, pageFaults uint64, err error) {
	file, err := os.Open("/proc/vmstat")
	if err != nil {
		return 0, 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		val, _ := strconv.ParseUint(fields[1], 10, 64)
		switch fields[0] {
		case "pgpgin":
			pageIn = val
		case "pgpgout":
			pageOut = val
		case "pgfault":
			pageFaults = val
		}
	}
	return pageIn, pageOut, pageFaults, nil
}

// getLinuxStats parses /proc/stat to get total interrupts and context switches
func getLinuxStats() (interrupts uint64, contextSwitches uint64, err error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "intr") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				interrupts, _ = strconv.ParseUint(fields[1], 10, 64)
			}
		} else if strings.HasPrefix(line, "ctxt") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				contextSwitches, _ = strconv.ParseUint(fields[1], 10, 64)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	return interrupts, contextSwitches, nil

}

func GetLinuxDiskPartitionsWithIO() ([]models.DiskPartition, error) {
	partitions, err := utils.GetDiskPartitions(false)
	if err != nil {
		return nil, err
	}

	// Capture initial diskstats snapshot
	initialStats, err := parseDiskStats()
	if err != nil {
		return nil, err
	}

	time.Sleep(1 * time.Second) // measure bytes/sec
	finalStats, err := parseDiskStats()
	if err != nil {
		return nil, err
	}

	var diskPartitions []models.DiskPartition
	for _, p := range partitions {
		UsedPercent, TotalDisk, DiskUsed, Fstype, err := utils.GetDiskUsage(p.Mountpoint)
		if err != nil {
			continue
		}

		// Health
		health := "healthy"
		if UsedPercent > 90 {
			health = "critical"
		} else if UsedPercent > 80 {
			health = "warning"
		}

		// Find matching device (strip /dev/)
		devName := strings.TrimPrefix(p.Device, "/dev/")
		readBytesSec := uint64(0)
		writeBytesSec := uint64(0)

		if s1, ok := initialStats[devName]; ok {
			if s2, ok := finalStats[devName]; ok {
				readBytesSec = (s2.readSectors - s1.readSectors) * 512
				writeBytesSec = (s2.writeSectors - s1.writeSectors) * 512
			}
		}

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

type diskStat struct {
	readSectors  uint64
	writeSectors uint64
}

// parse /proc/diskstats to extract read/write sectors
func parseDiskStats() (map[string]diskStat, error) {
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := make(map[string]diskStat)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}
		devName := fields[2]
		readSectors, _ := strconv.ParseUint(fields[5], 10, 64)
		writeSectors, _ := strconv.ParseUint(fields[9], 10, 64)
		stats[devName] = diskStat{readSectors, writeSectors}
	}
	return stats, scanner.Err()
}

func getOverallDiskStats() (readMBps, writeMBps float64, readIOPS, writeIOPS, totalIOPS uint64, busy, idle float64, err error) {
	// Take first snapshot
	snap1, err := parseDiskStats()
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, err
	}
	time.Sleep(1 * time.Second)
	snap2, err := parseDiskStats()
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, err
	}

	var totalReadSectors, totalWriteSectors uint64
	for dev, s1 := range snap1 {
		if s2, ok := snap2[dev]; ok {
			totalReadSectors += (s2.readSectors - s1.readSectors)
			totalWriteSectors += (s2.writeSectors - s1.writeSectors)
		}
	}

	// Convert sectors → MB (512 bytes per sector)
	readBytes := totalReadSectors * 512
	writeBytes := totalWriteSectors * 512
	readMBps = float64(readBytes) / (1024 * 1024)
	writeMBps = float64(writeBytes) / (1024 * 1024)

	// Approximate IOPS (1 sector = 1 I/O)
	readIOPS = totalReadSectors
	writeIOPS = totalWriteSectors
	totalIOPS = readIOPS + writeIOPS

	// Disk busy vs idle heuristic (simple: if IOPS > 0, disk busy)
	if totalIOPS > 0 {
		busy = 100
		idle = 0
	} else {
		busy = 0
		idle = 100
	}

	return readMBps, writeMBps, readIOPS, writeIOPS, totalIOPS, busy, idle, nil
}
