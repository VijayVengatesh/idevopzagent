//go:build windows
// +build windows

package systeminfo

import (
	"bytes"
	"fmt"
	"iDevopzAgent/internal/utils"
	"iDevopzAgent/models"
	"net"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"time"

	psnet "github.com/shirou/gopsutil/v3/net"
)

type WindowsCollector struct{}

func formatDuration(seconds uint64) string {
	d := time.Duration(seconds) * time.Second
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	mins := d / time.Minute
	d -= mins * time.Minute
	secs := d / time.Second
	return fmt.Sprintf("%d day(s) %d hr(s) %d min(s) %d sec(s)", days, hours, mins, secs)
}

func (w WindowsCollector) GetSystemSummary(userID string, machineId string) (*models.Systeminfo, error) {
	hostInfo, err := utils.HostInfo()
	if err != nil {
		return nil, err
	}

	cpuInfo, err := utils.GetCPUInfo()
	if err != nil {
		return nil, err
	}

	_, memTotal, _, _ := utils.GetMemoryUsage()
	diskInfo, _ := utils.GetDiskPartitions(false)
	netInterfaces, _ := utils.GetNetworkInterfaces()
	procsCount, _ := utils.GetProcessCount()
	currentUser, _ := user.Current()

	systemLogErrorCount, err := getWindowsEventLogErrorCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get system error log count: %v", err)
	}

	loginCount := getLoggedInUserCount()

	portCount, err := getOpenPortCount()
	if err != nil {
		fmt.Println("Error getting open port count:", err)
	} else {
		fmt.Println("Open Port Count:", portCount)
	}

	ip := getIP()

	summary := &models.Systeminfo{
		UserID:            userID,
		MachineID:         machineId,
		Hostname:          hostInfo.Hostname,
		IPAddress:         ip,
		OS:                fmt.Sprintf("%s %s (%s)", hostInfo.Platform, hostInfo.PlatformVersion, runtime.GOARCH),
		CPUModel:          cpuInfo[0].ModelName,
		CPUCores:          runtime.NumCPU(),
		RAMMB:             float64(memTotal) / (1024 * 1024),
		DiskCount:         len(diskInfo),
		SysLogsErrorCount: systemLogErrorCount,
		LoginCount:        loginCount,
		OpenPortCount:     portCount,
		Uptime:            formatDuration(hostInfo.Uptime),
		BootTime:          hostInfo.BootTime,
		TotalProcesses:    procsCount,
		NICCount:          len(netInterfaces),
		CurrentUser:       currentUser.Username,
	}

	return summary, nil
}

func getIP() string {
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		// Check interface is up and not loopback
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					ip := v.IP
					if ip.To4() != nil && !ip.IsLoopback() {
						return ip.String() // Return the first IPv4
					}
				case *net.IPAddr:
					ip := v.IP
					if ip.To4() != nil && !ip.IsLoopback() {
						return ip.String()
					}
				}
			}
		}
	}
	return "N/A" // Return first usable IP
}

func getWindowsEventLogErrorCount() (int, error) {
	cmd := exec.Command("powershell", "-Command", `
		(Get-WinEvent -LogName System | Where-Object { $_.LevelDisplayName -eq 'Error' } | Measure-Object).Count
	`)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("failed to execute command: %v", err)
	}

	countStr := strings.TrimSpace(out.String())
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse count: %v", err)
	}

	return count, nil
}

func getLoggedInUserCount() int {
	users, err := utils.GetLoggedInUsers()
	if err != nil {
		return 0
	}
	return len(users)
}

func getOpenPortCount() (int, error) {
	conns, err := psnet.Connections("all") // "tcp", "udp", or "all"
	if err != nil {
		return 0, err
	}

	count := 0
	for _, conn := range conns {
		if conn.Status == "LISTEN" || conn.Status == "ESTABLISHED" {
			count++
		}
	}
	return count, nil
}

func GetSystemInfoCollector() Collector {
	return WindowsCollector{}
}
