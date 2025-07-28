// internal/utilization/common.go
package utilization

import "iDevopzAgent/models"

type Collector interface {
	CpuUtilization(userID string) (*models.CpuUtilization, error)
	MemoryUtilization(userID string) (*models.MemoryUtilization, error)
	DiskUtilization(userID string) (*models.DiskUtilization, error)
}
