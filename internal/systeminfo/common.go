// internal/systeminfo/common.go
package systeminfo

import "iDevopzAgent/models"

type Collector interface {
	GetSystemSummary(userID string) (*models.Systeminfo, error)
}
