// internal/healthreport/common.go
package healthreport

import "iDevopzAgent/models"

type Collector interface {
	GenerateHealthReport(userId string) (*models.HealthReport, error)
}
