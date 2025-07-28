//go:build linux
// +build linux

package processdetails

import (
	"fmt"
	"iDevopzAgent/internal/utils"
	"iDevopzAgent/models"
	"os"
	"time"
)

type LinuxCollector struct{}

// Get all process details
func (l LinuxCollector) ListAllProcesses(userID string) ([]*models.ProcessInfo, error) {
	hostname, err := utils.GetHostName()
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %w", err)
	}

	procs, err := utils.GetAllProcesses()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	var results []*models.ProcessInfo
	for _, p := range procs {
		name, _ := p.Name()
		username, _ := p.Username()
		cpuPct, _ := p.CPUPercent()
		memPct, _ := p.MemoryPercent()
		threads, _ := p.NumThreads()
		priority, _ := p.Nice()

		var handles uint32

		fdPath := fmt.Sprintf("/proc/%d/fd", p.Pid)
		if fds, err := os.ReadDir(fdPath); err == nil {
			handles = uint32(len(fds))
		}

		processInfo := &models.ProcessInfo{
			UserID:        userID,
			Hostname:      hostname,
			PID:           p.Pid,
			Name:          name,
			Username:      username,
			CPUPercent:    cpuPct,
			MemoryPercent: memPct,
			ThreadCount:   threads,
			HandleCount:   handles,
			Priority:      string(priority),
			Timestamp:     time.Now().Unix(),
		}
		results = append(results, processInfo)
	}

	return results, nil
}

func GetProcessCollector() Collector {
	return LinuxCollector{}
}
