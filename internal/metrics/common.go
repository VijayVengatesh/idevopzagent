// internal/metrics/common.go
package metrics

import "iDevopzAgent/models"

type Metrics struct {
	CPUUsage    float64
	MemoryUsage float64
}

type Collector interface {
	Collect() Metrics
	MetricsCollect(userID string) (*models.Metrics, error)
}
