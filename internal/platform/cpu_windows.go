//go:build windows
// +build windows

package platform

type WindowsMetrics struct{}

func (w WindowsMetrics) GetCPUUsage() string {
	return "Windows CPU Usage"
}
