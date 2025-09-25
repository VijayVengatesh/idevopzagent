//go:build windows
// +build windows

package processdetails

import (
	"fmt"
	"iDevopzAgent/internal/utils"
	"iDevopzAgent/models"
	"os"
	"runtime"
	"time"

	"unsafe"

	"golang.org/x/sys/windows"
)

type WindowsCollector struct{}

var (
	modNtDll                      = windows.NewLazySystemDLL("ntdll.dll")
	procNtQueryInformationProcess = modNtDll.NewProc("NtQueryInformationProcess")
)

const (
	ProcessHandleCount = 51
)

// Priority class mapping for Windows
func getWindowsPriorityClass(priority int32) string {
	switch {
	case priority <= 4:
		return "REALTIME_PRIORITY"
	case priority <= 8:
		return "HIGH_PRIORITY"
	case priority <= 13:
		return "ABOVE_NORMAL_PRIORITY"
	case priority <= 18:
		return "NORMAL_PRIORITY"
	case priority <= 23:
		return "BELOW_NORMAL_PRIORITY"
	default:
		return "IDLE_PRIORITY"
	}
}

// Get all process details
func (w WindowsCollector) ListAllProcesses(userID string, machineId string) ([]*models.ProcessInfo, error) {
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
		path, err := p.Exe()
		if err != nil {
			// Skip system processes that don't have accessible paths (like PID 0)
			if p.Pid == 0 || p.Pid == 4 {
				path = "[System Process]"
			} else {
				continue // Skip other inaccessible processes
			}
		}
		var handles uint32

		handleCount, err := GetWindowsHandleCount(p.Pid)
		if err == nil {
			handles = handleCount
		}

		priorityLabel := fmt.Sprintf("Nice %d", priority)
		if runtime.GOOS == "windows" {
			priorityLabel = getWindowsPriorityClass(priority)
		}

		processInfo := &models.ProcessInfo{
			UserID:        userID,
			MachineID:     machineId,
			Path:          path,
			Hostname:      hostname,
			PID:           p.Pid,
			Name:          name,
			Username:      username,
			CPUPercent:    cpuPct,
			MemoryPercent: memPct,
			ThreadCount:   threads,
			HandleCount:   handles,
			Priority:      priorityLabel,
			Timestamp:     time.Now().Unix(),
		}
		results = append(results, processInfo)
	}

	return results, nil
}

//top 5 cpu process

func (w WindowsCollector) ListTop5CpuProcess(userID string, machineId string) ([]*models.Process, error) {
	topCpuProcesses, err := utils.GetTopProcessesByCPU(5)
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	var result []*models.Process
	for _, proc := range topCpuProcesses {
		result = append(result, &models.Process{
			UserID:    userID,
			MachineID: machineId,
			Hostname:  hostname,
			PID:       fmt.Sprintf("%d", proc.PID),
			Usage:     proc.CPUPercent,
			Command:   proc.Name,
		})
	}

	return result, nil
}

func (w WindowsCollector) ListTop5MemoryProcess(userID string, machineId string) ([]*models.Process, error) {
	topCpuProcesses, err := utils.GetTopProcessesByMemory(5)
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	var result []*models.Process
	for _, proc := range topCpuProcesses {
		result = append(result, &models.Process{
			UserID:    userID,
			MachineID: machineId,
			Hostname:  hostname,
			PID:       fmt.Sprintf("%d", proc.PID),
			Usage:     float64(proc.MemPercent),
			Command:   proc.Name,
		})
	}

	return result, nil
}
func GetWindowsHandleCount(pid int32) (uint32, error) {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(handle)

	var handleCount uint32
	ret, _, err := procNtQueryInformationProcess.Call(
		uintptr(handle),
		uintptr(ProcessHandleCount),
		uintptr(unsafe.Pointer(&handleCount)),
		unsafe.Sizeof(handleCount),
		0,
	)
	if ret != 0 {
		return 0, err
	}
	return handleCount, nil
}

func GetProcessCollector() Collector {
	return WindowsCollector{}
}
