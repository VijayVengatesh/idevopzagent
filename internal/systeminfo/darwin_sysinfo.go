//go:build darwin

package systeminfo

import (
	"fmt"
	"iDevopzAgent/models"
	"net"
	"os/user"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	psnet "github.com/shirou/gopsutil/v3/net"
)

type DarwinSystemInfoCollector struct{}

func GetSystemInfoCollector() Collector {
	return &DarwinSystemInfoCollector{}
}

func (s *DarwinSystemInfoCollector) GetSystemSummary(userID, machineID string) (*models.Systeminfo, error) {
	// Host information
	hostInfo, _ := host.Info()

	// CPU information
	cpuInfo, _ := cpu.Info()
	cpuCount, _ := cpu.Counts(true)

	// Memory information
	memInfo, _ := mem.VirtualMemory()

	// Disk information
	diskPartitions, _ := disk.Partitions(false)
	netInterfaces, _ := net.Interfaces()
	procsCount := getProcessCount()
	currentUser := getCurrentUser()

	errorLogCount := 0 // macOS system logs would need different approach
	loginCount := getLoggedInUserCount()
	openPortCount := getOpenPortCount()
	ip := getIP()

	var cpuModel string
	if len(cpuInfo) > 0 {
		cpuModel = cpuInfo[0].ModelName
	}

	systemInfo := &models.Systeminfo{
		UserID:            userID,
		MachineID:         machineID,
		Hostname:          hostInfo.Hostname,
		IPAddress:         ip,
		OS:                fmt.Sprintf("%s %s (%s)", hostInfo.Platform, hostInfo.PlatformVersion, hostInfo.KernelArch),
		CPUModel:          cpuModel,
		CPUCores:          cpuCount,
		RAMMB:             float64(memInfo.Total) / (1024 * 1024),
		DiskCount:         len(diskPartitions),
		SysLogsErrorCount: errorLogCount,
		LoginCount:        loginCount,
		OpenPortCount:     openPortCount,
		Uptime:            formatDuration(hostInfo.Uptime),
		BootTime:          hostInfo.BootTime,
		TotalProcesses:    procsCount,
		NICCount:          len(netInterfaces),
		CurrentUser:       currentUser,
	}

	return systemInfo, nil
}

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

func getIP() string {
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
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
						return ip.String()
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
	return "N/A"
}

func getLoggedInUserCount() int {
	// macOS specific implementation would be needed
	return 1 // Default to 1 for current user
}

func getOpenPortCount() int {
	conns, err := psnet.Connections("all")
	if err != nil {
		return 0
	}
	count := 0
	for _, conn := range conns {
		if conn.Status == "LISTEN" || conn.Status == "ESTABLISHED" {
			count++
		}
	}
	return count
}

func getProcessCount() int {
	procs, err := process.Processes()
	if err != nil {
		return 0
	}
	return len(procs)
}

func getCurrentUser() string {
	user, err := user.Current()
	if err != nil {
		return "unknown"
	}
	return user.Username
}