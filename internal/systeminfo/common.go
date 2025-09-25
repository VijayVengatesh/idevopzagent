// internal/systeminfo/common.go
package systeminfo

import "iDevopzAgent/models"

type Collector interface {
	GetSystemSummary(userID string, machineID string) (*models.Systeminfo, error)
}
