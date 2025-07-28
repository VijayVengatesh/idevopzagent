package models

type ProcessInfo struct {
	UserID        string  `json:"user_id"`
	Hostname      string  `json:"hostname"`
	PID           int32   `json:"pid"`
	Name          string  `json:"name"`
	Username      string  `json:"user_name"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float32 `json:"memory_percent"`
	ThreadCount   int32   `json:"thread_count"`
	HandleCount   uint32  `json:"handle_count"`
	Priority      string  `json:"priority"`
	Timestamp     int64   `json:"timestamp"`
}
