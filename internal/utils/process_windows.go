//go:build windows

package utils

import (
	"os/exec"
	"strconv"
	"strings"
)

// IsProcessRunning checks if a process with given PID and name is running using Windows tasklist
func IsProcessRunning(pid int, name string) bool {
	cmd := exec.Command("tasklist", "/FI", "PID eq "+strconv.Itoa(pid), "/FO", "CSV", "/NH")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		// Parse CSV format: "ProcessName","PID","SessionName","Session#","MemUsage"
		fields := strings.Split(line, ",")
		if len(fields) >= 2 {
			processName := strings.Trim(fields[0], "\"")
			processPID := strings.Trim(fields[1], "\"")
			
			if processPID == strconv.Itoa(pid) && processName == name {
				return true
			}
		}
	}
	
	return false
}