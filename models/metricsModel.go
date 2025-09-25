package models

type DiskPartition struct {
	Device        string  `json:"device"`
	Mountpoint    string  `json:"mountpoint"`
	Fstype        string  `json:"fstype"`
	Total         uint64  `json:"total"`
	Used          uint64  `json:"used"`
	UsedPercent   float64 `json:"used_percent"`
	Healthy       string  `json:"healthy"`
	ReadBytesSec  uint64  `json:"read_bytes_sec"`
	WriteBytesSec uint64  `json:"write_bytes_sec"`
}
type Metrics struct {
	UserID               string          `json:"user_id"`
	Hostname             string          `json:"hostname"`
	MachineID            string          `json:"machineId"`
	CPUPercent           float64         `json:"cpu_percent"`
	CPUPerCore           []float64       `json:"cpu_per_core"` // per-core usage
	SystemIdle           float64         `json:"system_idle"`
	MemoryUsed           uint64          `json:"memory_used"`
	MemoryTotal          uint64          `json:"memory_total"`
	MemoryFree           uint64          `json:"memory_free"`
	MemoryPercent        float64         `json:"memory_percent"`
	SwapMemoryUsed       uint64          `json:"swap_memory_used"`
	SwapMemoryTotal      uint64          `json:"swap_memory_total"`
	SwapMemoryFree       uint64          `json:"swap_memory_free"`
	SwapMemUsagePercent  float64         `json:"swap_mem_usage_percent"`
	DiskUsed             uint64          `json:"disk_used"`
	DiskTotal            uint64          `json:"disk_total"`
	Uptime               uint64          `json:"uptime_seconds"`
	MetricGetTime        string          `json:"metric_get_time"`
	Status               string          `json:"status"` // up, down, trouble, critical
	Timestamp            int64           `json:"timestamp"`
	Os                   string          `json:"os`
	Interrupts           uint64          `json:"interrupts"`
	ContextSwitches      uint64          `json:"context_switches"`
	PagesReads           uint            `json:"pages_reads"`
	PagesWrites          uint            `json:"pages_writes"`
	PageFaults           uint            `json:"page_faults"`
	DiskPartitions       []DiskPartition `json:"disk_partitions"`
	OverallDiskReadMBps  float64         `json:"overall_disk_read_mb"`
	OverallDiskWriteMBps float64         `json:"overall_disk_write_mb"`
	OverallDiskReadIOPS  uint64          `json:"overall_disk_read_iops"`
	OverallDiskWriteIOPS uint64          `json:"overall_disk_write_iops"`
	OverallDiskIOPS      uint64          `json:"overall_disk_iops"`
	OverallDiskIdle      float64         `json:"overall_disk_idle_percent"`
	OverallDiskBusy      float64         `json:"overall_disk_busy_percent"`

	// up, down, trouble, critical

}
