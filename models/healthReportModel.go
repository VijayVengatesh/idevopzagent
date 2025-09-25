package models

type HealthReport struct {
	UserID        string  `json:"user_id"`
	MachineID     string  `json:"machineId"`
	Hostname      string  `json:"hostname"`
	Availability  string  `json:"availability"`
	CPUPercent    float64 `json:"cpu"`
	MemoryPercent float64 `json:"memory"`
	DiskPercent   float64 `json:"disk"`
	Downtimes     int     `json:"downtimes"`
	MetricGetTime string  `json:"metric_get_time"`

	SLA string `json:"sla_achieved"`
}
