package utils

import (
	"runtime"

	"github.com/shirou/gopsutil/v3/load"
)

// IsLoadSupported checks if load averages are supported on this OS
func IsLoadSupported() bool {
	return runtime.GOOS != "windows"
}

// GetLoadAverage returns system load average (1m, 5m, 15m)
func GetLoadAverage() (*load.AvgStat, error) {
	if !IsLoadSupported() {
		return nil, nil // or return custom error
	}
	return load.Avg()
}

// GetLoadMisc returns load-related misc stats (Linux only):
// - ProcsRunning: number of processes currently running
// - ProcsBlocked: number of processes blocked waiting for I/O
func GetLoadMisc() (*load.MiscStat, error) {
	if !IsLoadSupported() {
		return nil, nil
	}
	return load.Misc()
}
