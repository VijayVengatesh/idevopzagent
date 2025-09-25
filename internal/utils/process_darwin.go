//go:build darwin

package utils

import (
	"os/exec"
	"strconv"
	"strings"
)

// IsProcessRunning checks if a process with given PID and name is running using macOS ps command
func IsProcessRunning(pid int, name string) bool {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	processName := strings.TrimSpace(string(output))
	return processName == name
}