//go:build linux
// +build linux

package platform

type LinuxMetrics struct{}

func (l LinuxMetrics) GetCPUUsage() string {
	return "Linux CPU Usage"
}
