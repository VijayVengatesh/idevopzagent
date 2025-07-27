package platform

type MetricsCollector interface {
	GetCPUUsage() string
	GetMemoryUsage() string
}

func GetMetricsCollector() MetricsCollector {
	return getMetricsCollector()
}
