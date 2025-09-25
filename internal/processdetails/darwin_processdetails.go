//go:build darwin

package processdetails

import (
	"bufio"
	"fmt"
	"iDevopzAgent/models"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type DarwinProcessCollector struct{}

func GetProcessCollector() Collector {
	return &DarwinProcessCollector{}
}

func (p *DarwinProcessCollector) ListAllProcesses(userID, machineID string) ([]*models.ProcessInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var processList []*models.ProcessInfo
	for _, proc := range procs {
		name, _ := proc.Name()
		cpuPercent, _ := proc.CPUPercent()
		memPercent, _ := proc.MemoryPercent()
		username, _ := proc.Username()
		threads, _ := proc.NumThreads()
		nice, _ := proc.Nice()

		// Get file descriptors count (handles equivalent)
		var handles uint32
		if fds, err := proc.OpenFiles(); err == nil {
			handles = uint32(len(fds))
		}

		processInfo := &models.ProcessInfo{
			UserID:        userID,
			MachineID:     machineID,
			Hostname:      hostname,
			PID:           proc.Pid,
			Name:          name,
			Username:      username,
			CPUPercent:    cpuPercent,
			MemoryPercent: float32(memPercent),
			ThreadCount:   threads,
			HandleCount:   handles,
			Priority:      int(nice),
			Timestamp:     time.Now().Unix(),
		}

		processList = append(processList, processInfo)
	}

	return processList, nil
}

func (p *DarwinProcessCollector) ListTop5CpuProcess(userID, machineID string) ([]*models.Process, error) {
	cmd := exec.Command("ps", "-eo", "pid,comm,%cpu", "-r")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running ps command for top5 cpu process: %w", err)
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
		if lineNum == 1 || lineNum > 6 {
			continue // skip header and limit to top 5
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
			MachineID: machineID,
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

func (p *DarwinProcessCollector) ListTop5MemoryProcess(userID, machineID string) ([]*models.Process, error) {
	cmd := exec.Command("ps", "-eo", "pid,comm,%mem", "-m")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running ps command for top5 memory process: %w", err)
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
		if lineNum == 1 || lineNum > 6 {
			continue // skip header and limit to top 5
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		pid := fields[0]
		command := fields[1]
		memStr := fields[2]

		memPercent, err := strconv.ParseFloat(memStr, 64)
		if err != nil {
			memPercent = 0.0
		}

		processes = append(processes, &models.Process{
			UserID:    userID,
			MachineID: machineID,
			Hostname:  hostname,
			PID:       pid,
			Usage:     memPercent,
			Command:   command,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return processes, nil
}
