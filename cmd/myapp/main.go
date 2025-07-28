// cmd/yourapp/main.go
package main

import (
	"fmt"
	"iDevopzAgent/configs"
	"iDevopzAgent/internal/metrics"
	"iDevopzAgent/internal/systeminfo"
)

func main() {

	userID, err := configs.LoadUserID()
	if err != nil || userID == "" {
		userID = configs.PromptAndSaveUserID()
	}
	// Metrics
	collector := metrics.GetCollector()
	m := collector.Collect()
	y, err := collector.MetricsCollect(userID)
	fmt.Println("collected metrics", y, err)
	// System Info
	fetcher := systeminfo.GetFetcher()
	info := fetcher.GetInfo()

	fmt.Println("System:", info.Platform)
	fmt.Println("Hostname:", info.Hostname)
	fmt.Println("CPU Usage:", m.CPUUsage)
	fmt.Println("Memory Usage:", m.MemoryUsage)

}
