// internal/utilization/common.go
package utilization

import "iDevopzAgent/models"

type Collector interface {
	CpuUtilization(userID string, machineID string) (*models.CpuUtilization, error)
	MemoryUtilization(userID string, machineID string) (*models.MemoryUtilization, error)
	DiskUtilization(userID string, machineID string) (*models.DiskUtilization, error)
}
