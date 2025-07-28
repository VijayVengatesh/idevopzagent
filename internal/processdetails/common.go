// internal/processdetails/common.go
package processdetails

import "iDevopzAgent/models"

type Collector interface {
	ListAllProcesses(userID string) ([]*models.ProcessInfo, error)
}
