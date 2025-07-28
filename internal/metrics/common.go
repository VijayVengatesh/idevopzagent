// internal/metrics/common.go
package metrics

import "iDevopzAgent/models"

type Collector interface {
	MetricsCollect(userID string) (*models.Metrics, error)
}
