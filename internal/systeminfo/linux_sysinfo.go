//go:build linux
// +build linux

package systeminfo

import (
	"fmt"
	"iDevopzAgent/internal/utils"
	"iDevopzAgent/models"
	"net"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	psnet "github.com/shirou/gopsutil/v3/net" // alias to avoid std lib conflict
)

type LinuxCollector struct{}

func (l LinuxCollector) GetSystemSummary(userID string) (*models.Systeminfo, error) {
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

	errorLogCount := countSyslogErrors()
	loginCount := getLoggedInUserCount()
	openPortcount, err := GetOpenPortCount()
	if err != nil {
		fmt.Println("Error:", err)

	}

	ip := getIP()

	summary := &models.Systeminfo{
		UserID:            userID,
		Hostname:          hostInfo.Hostname,
		IPAddress:         ip,
		OS:                fmt.Sprintf("%s %s (%s)", hostInfo.Platform, hostInfo.PlatformVersion, runtime.GOARCH),
		CPUModel:          cpuInfo[0].ModelName,
		CPUCores:          runtime.NumCPU(),
		RAMMB:             float64(memTotal) / (1024 * 1024),
		DiskCount:         len(diskInfo),
		SysLogsErrorCount: errorLogCount,
		LoginCount:        loginCount,
		OpenPortCount:     openPortcount,
		Uptime:            formatDuration(hostInfo.Uptime),
		BootTime:          hostInfo.BootTime,
		TotalProcesses:    procsCount,
		NICCount:          len(netInterfaces),
		CurrentUser:       currentUser.Username,
	}

	return summary, nil
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
func countSyslogErrors() int {
	content, err := os.ReadFile("/var/log/syslog")
	if err != nil {
		return 0
	}
	count := 0
	for _, line := range strings.Split(string(content), "\n") {
		if strings.Contains(line, "error") || strings.Contains(line, "ERROR") {
			count++
		}
	}
	return count
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
func getLoggedInUserCount() int {
	users, err := utils.GetLoggedInUsers()
	if err != nil {
		return 0
	}
	return len(users)
}

func GetOpenPortCount() (int, error) {
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
	return LinuxCollector{}
}
