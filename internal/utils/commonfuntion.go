package utils

import (
	"os"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetNumCPU() int {
	return runtime.NumCPU()
}
func GetOS() string {
	return runtime.GOOS
}
func GetHostName() (string, error) {
	return os.Hostname()
}

func GetCPUPercentage() (float64, error) {
	percentages, err := cpu.Percent(0, false)
	if err != nil || len(percentages) == 0 {
		return 0, err
	}
	return percentages[0], nil
}

func GetMemoryUsage() (usedPercent float64, total uint64, used uint64, err error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, 0, err
	}
	return vmStat.UsedPercent, vmStat.Total, vmStat.Used, nil
}

func GetDiskUsage(path string) (usedPercent float64, total uint64, used uint64, err error) {
	diskStat, err := disk.Usage(path)
	if err != nil {
		return 0, 0, 0, err
	}
	return diskStat.UsedPercent, diskStat.Total, diskStat.Used, nil
}

func GetUptime() (uint64, error) {
	info, err := host.Info()
	if err != nil {
		return 0, err
	}
	return info.Uptime, nil
}
