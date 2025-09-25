package models

type ProcessInfo struct {
	UserID        string  `json:"user_id"`
	MachineID     string  `json:"machineId"`
	Hostname      string  `json:"hostname"`
	Path          string  `json:"path"`
	PID           int32   `json:"pid"`
	Name          string  `json:"name"`
	Username      string  `json:"user_name"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float32 `json:"memory_percent"`
	ThreadCount   int32   `json:"thread_count"`
	HandleCount   uint32  `json:"handle_count"`
	Priority      any     `json:"priority"`
	Timestamp     int64   `json:"timestamp"`
}

type Process struct {
	UserID    string  `json:"user_id"`
	MachineID string  `json:"machineId"`
	Hostname  string  `json:"hostname"`
	PID       string  `json:"pid"`
	Usage     float64 `json:"usage"`
	Command   string  `json:"command"`
}
