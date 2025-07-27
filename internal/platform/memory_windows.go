//go:build windows
// +build windows

package platform

func (w WindowsMetrics) GetMemoryUsage() string {
	return "Windows Memory Usage"
}

func getMetricsCollector() MetricsCollector {
	return WindowsMetrics{}
}
