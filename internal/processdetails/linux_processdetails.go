//go:build linux
// +build linux

package processdetails

import (
	"bufio"
	"fmt"
	"iDevopzAgent/internal/utils"
	"iDevopzAgent/models"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type LinuxCollector struct{}

// Get all process details
func (l LinuxCollector) ListAllProcesses(userID string, machineId string) ([]*models.ProcessInfo, error) {
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
			MachineID:     machineId,
			Hostname:      hostname,
			PID:           p.Pid,
			Name:          name,
			Username:      username,
			CPUPercent:    cpuPct,
			MemoryPercent: memPct,
			ThreadCount:   threads,
			HandleCount:   handles,
			Priority:      int(priority),
			Timestamp:     time.Now().Unix(),
		}
		results = append(results, processInfo)
	}

	return results, nil
}

//top 5 cpu process

func (l LinuxCollector) ListTop5CpuProcess(userID string, machineId string) ([]*models.Process, error) {
	cmd := exec.Command("bash", "-c", "ps -eo pid,comm,%cpu --sort=-%cpu | head -n 6")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running ps command while top5 cpu process: %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var processes []*models.Process
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		if lineNum == 1 {
			continue // skip header
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		pid := fields[0]
		command := fields[1]
		cpuStr := fields[2]

		cpuPercent, err := strconv.ParseFloat(cpuStr, 64)
		if err != nil {
			cpuPercent = 0.0
		}

		processes = append(processes, &models.Process{
			UserID:    userID,
			MachineID: machineId,
			Hostname:  hostname,
			PID:       pid,
			Usage:     cpuPercent,
			Command:   command,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return processes, nil
}

func (l LinuxCollector) ListTop5MemoryProcess(userID string, machineId string) ([]*models.Process, error) {
	cmd := exec.Command("bash", "-c", "ps -eo pid,comm,%mem --sort=-%mem | head -n 6")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running ps command while top5 mem process : %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var processes []*models.Process
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		if lineNum == 1 {
			continue // skip header
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		pid := fields[0]
		command := fields[1]
		cpuStr := fields[2]

		cpuPercent, err := strconv.ParseFloat(cpuStr, 64)
		if err != nil {
			cpuPercent = 0.0
		}

		processes = append(processes, &models.Process{
			UserID:   userID,
			Hostname: hostname,
			PID:      pid,
			Usage:    cpuPercent,
			Command:  command,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return processes, nil
}

func GetProcessCollector() Collector {
	return LinuxCollector{}
}
