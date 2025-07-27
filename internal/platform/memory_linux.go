//go:build linux
// +build linux

package platform

func (l LinuxMetrics) GetMemoryUsage() string {
	return "Linux Memory Usage"
}

func getMetricsCollector() MetricsCollector {
	return LinuxMetrics{}
}
