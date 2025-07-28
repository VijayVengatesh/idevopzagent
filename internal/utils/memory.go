package utils

import (
	"github.com/shirou/gopsutil/v3/mem"
)

// GetMemoryUsage returns used percentage, total memory, and used memory
func GetMemoryUsage() (usedPercent float64, total uint64, used uint64, err error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, 0, err
	}
	return vmStat.UsedPercent, vmStat.Total, vmStat.Used, nil
}

// GetSwapUsage returns swap memory stats: used %, total, used
func GetSwapUsage() (usedPercent float64, total uint64, used uint64, err error) {
	swap, err := mem.SwapMemory()
	if err != nil {
		return 0, 0, 0, err
	}
	return swap.UsedPercent, swap.Total, swap.Used, nil
}
