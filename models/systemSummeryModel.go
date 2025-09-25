package models

type Systeminfo struct {
	UserID            string  `json:"user_id"`
	MachineID         string  `json:"machineId"`
	Hostname          string  `json:"hostname"`
	IPAddress         string  `json:"ip_address"`
	OS                string  `json:"os"`
	CPUModel          string  `json:"cpu_model"`
	CPUCores          int     `json:"cpu_cores"`
	RAMMB             float64 `json:"ram_mb"`
	DiskCount         int     `json:"disk_count"`
	SysLogsErrorCount int     `json:"sys_logs_errors"`
	Uptime            string  `json:"uptime"`
	BootTime          uint64  `json:"boot_time"`
	TotalProcesses    int     `json:"total_processes"`
	NICCount          int     `json:"nic_count"`
	LoginCount        int     `json:"login_count"`
	OpenPortCount     int     `json:"open_port_count"`
	CurrentUser       string  `json:"current_user"`
}
