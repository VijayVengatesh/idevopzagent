package models

type CpuUtilization struct {
	UserID     string  `json:"user_id"`
	MachineID  string  `json:"machine_id"`
	Hostname   string  `json:"hostname"`
	CPUPercent float64 `json:"cpu_percent"`
}

type MemoryUtilization struct {
	UserID        string  `json:"user_id"`
	MachineID     string  `json:"machine_id"`
	Hostname      string  `json:"hostname"`
	MemoryPercent float64 `json:"memory_percent"`
}
type DiskUtilization struct {
	UserID      string  `json:"user_id"`
	MachineID   string  `json:"machine_id"`
	Hostname    string  `json:"hostname"`
	UsedPercent float64 `json:"used_percent"`
	FreePercent float64 `json:"disk_free_percent"`
}
