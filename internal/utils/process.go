package utils

import (
	"sort"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo holds basic info about a process
type ProcessInfo struct {
	PID        int32
	Name       string
	CPUPercent float64
	MemPercent float32
}

// GetAllProcesses returns a slice of all process objects
func GetAllProcesses() ([]*process.Process, error) {
	return process.Processes()
}

// GetTopProcessesByCPU returns top N processes by CPU usage
func GetTopProcessesByCPU(limit int) ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var infos []ProcessInfo
	for _, p := range procs {
		name, _ := p.Name()
		cpu, _ := p.CPUPercent()
		mem, _ := p.MemoryPercent()

		if cpu > 0 { // Skip idle processes
			infos = append(infos, ProcessInfo{
				PID:        p.Pid,
				Name:       name,
				CPUPercent: cpu,
				MemPercent: mem,
			})
		}
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].CPUPercent > infos[j].CPUPercent
	})

	if len(infos) > limit {
		infos = infos[:limit]
	}
	return infos, nil
}

// GetTopProcessesByMemory returns top N processes by memory usage
func GetTopProcessesByMemory(limit int) ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var infos []ProcessInfo
	for _, p := range procs {
		name, _ := p.Name()
		cpu, _ := p.CPUPercent()
		mem, _ := p.MemoryPercent()

		if mem > 0 {
			infos = append(infos, ProcessInfo{
				PID:        p.Pid,
				Name:       name,
				CPUPercent: cpu,
				MemPercent: mem,
			})
		}
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].MemPercent > infos[j].MemPercent
	})

	if len(infos) > limit {
		infos = infos[:limit]
	}
	return infos, nil
}

// GetProcessCount returns the number of running processes
func GetProcessCount() (int, error) {
	procs, err := process.Processes()
	if err != nil {
		return 0, err
	}
	return len(procs), nil
}


