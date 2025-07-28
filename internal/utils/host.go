package utils

import (
	"os"
	"runtime"

	"github.com/shirou/gopsutil/v3/host"
)

// GetHostName returns the hostname of the machine
func GetHostName() (string, error) {
	return os.Hostname()
}

// GetOS returns the OS type (windows/linux/darwin)
func GetOS() string {
	return runtime.GOOS
}

// GetUptime returns system uptime in seconds
func GetUptime() (uint64, error) {
	info, err := host.Info()
	if err != nil {
		return 0, err
	}
	return info.Uptime, nil
}

// HostInfo returns detailed host/system information
func HostInfo() (*host.InfoStat, error) {
	return host.Info()
}
