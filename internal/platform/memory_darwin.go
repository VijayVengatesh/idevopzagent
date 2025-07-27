//go:build darwin
// +build darwin

package platform

func (m MacMetrics) GetMemoryUsage() string {
	return "macOS Memory Usage"
}

func getMetricsCollector() MetricsCollector {
	return MacMetrics{}
}
