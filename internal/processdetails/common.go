// internal/processdetails/common.go
package processdetails

import "iDevopzAgent/models"

type Collector interface {
	ListAllProcesses(userID string, machineId string) ([]*models.ProcessInfo, error)
	ListTop5MemoryProcess(userID string, machineId string) ([]*models.Process, error)
	ListTop5CpuProcess(userID string, machineId string) ([]*models.Process, error)
}
