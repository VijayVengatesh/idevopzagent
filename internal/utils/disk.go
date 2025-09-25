package utils

import (
	"github.com/shirou/gopsutil/v3/disk"
)

// GetDiskUsage returns used percent, total and used disk space of the given path
func GetDiskUsage(path string) (usedPercent float64, total uint64, used uint64, fstype string, err error) {
	diskStat, err := disk.Usage(path)
	if err != nil {
		return 0, 0, 0, "", err
	}
	return diskStat.UsedPercent, diskStat.Total, diskStat.Used, diskStat.Fstype, nil
}

// GetDiskPartitions returns all mounted partitions
func GetDiskPartitions(all bool) ([]disk.PartitionStat, error) {
	return disk.Partitions(all)
}

// GetIOCounters returns disk I/O stats for all devices
func GetIOCounters() (map[string]disk.IOCountersStat, error) {
	return disk.IOCounters()
}
