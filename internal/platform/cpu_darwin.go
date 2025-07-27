//go:build darwin
// +build darwin

package platform

type MacMetrics struct{}

func (m MacMetrics) GetCPUUsage() string {
	return "macOS CPU Usage"
}
