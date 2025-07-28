package utils

import (
	"github.com/shirou/gopsutil/v3/net"
)

// NetIOInfo holds summarized network IO statistics
type NetIOInfo struct {
	Name        string
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
}

// GetNetworkInterfaces returns all network interfaces (physical & virtual)
func GetNetworkInterfaces() ([]net.InterfaceStat, error) {
	return net.Interfaces()
}

// GetNetworkIO returns a map of network IO stats per interface
func GetNetworkIO() (map[string]NetIOInfo, error) {
	ioCounters, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]NetIOInfo)
	for _, io := range ioCounters {
		stats[io.Name] = NetIOInfo{
			Name:        io.Name,
			BytesSent:   io.BytesSent,
			BytesRecv:   io.BytesRecv,
			PacketsSent: io.PacketsSent,
			PacketsRecv: io.PacketsRecv,
		}
	}
	return stats, nil
}

// GetTotalNetworkIO returns total network stats (all interfaces combined)
func GetTotalNetworkIO() (*NetIOInfo, error) {
	ioCounters, err := net.IOCounters(false)
	if err != nil || len(ioCounters) == 0 {
		return nil, err
	}

	total := ioCounters[0]
	return &NetIOInfo{
		Name:        "total",
		BytesSent:   total.BytesSent,
		BytesRecv:   total.BytesRecv,
		PacketsSent: total.PacketsSent,
		PacketsRecv: total.PacketsRecv,
	}, nil
}
