// internal/metrics/common.go
package metrics

import "iDevopzAgent/models"

type Collector interface {
	MetricsCollect(userID string, machineId string) (*models.Metrics, error)
}
