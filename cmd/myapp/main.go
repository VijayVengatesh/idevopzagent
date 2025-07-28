// cmd/yourapp/main.go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"iDevopzAgent/configs"
	"iDevopzAgent/internal/healthreport"
	"iDevopzAgent/internal/metrics"
	"iDevopzAgent/internal/processdetails"
	"iDevopzAgent/internal/systeminfo"
	"iDevopzAgent/internal/utilization"
	"iDevopzAgent/internal/utils"
	"iDevopzAgent/sender"
)

func main() {
	userID, err := configs.LoadUserID()
	if err != nil || userID == "" {
		userID = configs.PromptAndSaveUserID()
	}

	go collectMetrics(userID)
	go collectUtilization(userID)
	go collectHealthReport(userID)
	go collectProcessDetails(userID)
	go collectSystemInfo(userID)

	// Prevent the main function from exiting
	select {}
}

func collectMetrics(userID string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	collector := metrics.GetCollector()

	for {
		select {
		case <-ticker.C:
			y, err := collector.MetricsCollect(userID)
			sender.SendToMetricsAPI(y)
			fmt.Println("collected metrics", y, err)
		}
	}
}

func collectUtilization(userID string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	u := utilization.UtilizationCollector()

	for {
		select {
		case <-ticker.C:
			if cpuUtil, err := u.CpuUtilization(userID); err == nil {
				fmt.Println("cpu utilization", cpuUtil)
				sender.SendCpuUtilizationToAPI(cpuUtil)

			}
			if memUtil, err := u.MemoryUtilization(userID); err == nil {
				fmt.Println("memory utilization", memUtil)
				sender.SendMemmoryUtilizationToAPI(memUtil)

			}
			if diskUtil, err := u.DiskUtilization(userID); err == nil {
				fmt.Println("disk utilization", diskUtil)
				sender.SendDiskUtilizationToAPI(diskUtil)

			}
		}
	}
}

func collectHealthReport(userID string) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	h := healthreport.GetHealthReportCollector()

	for {
		select {
		case <-ticker.C:
			health, err := h.GenerateHealthReport(userID)
			if err != nil {
				fmt.Println("Error collecting healthReport:", err)
			} else {
				fmt.Println("healthReport", health)
				sender.SendToHealthReportAPI(health)

			}
		}
	}
}

func collectProcessDetails(userID string) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	processUtil := processdetails.GetProcessCollector()

	for {
		select {
		case <-ticker.C:
			p, err := processUtil.ListAllProcesses(userID)
			if err == nil {
				data, _ := json.MarshalIndent(p, "", "  ")
				fmt.Println("process details", string(data), len(p))
				sender.SendProcessList(p)

			}

			if top5Cpu, err := processUtil.ListTop5CpuProcess(userID); err == nil {
				data, _ := json.MarshalIndent(top5Cpu, "", "  ")
				fmt.Println("Top 5 processes by CPU:\n", string(data))
				sender.Top5Cpu(top5Cpu)

			}

			if top5Mem, err := processUtil.ListTop5MemoryProcess(userID); err == nil {
				data, _ := json.MarshalIndent(top5Mem, "", "  ")
				fmt.Println("Top 5 processes by Memory:\n", string(data))
				sender.Top5Memory(top5Mem)

			}

			if count, err := utils.GetProcessCount(); err == nil {
				fmt.Println("Process count:", count)
			}
		}
	}
}

func collectSystemInfo(userID string) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	systemInfoCollector := systeminfo.GetSystemInfoCollector()

	for {
		select {
		case <-ticker.C:
			sys, err := systemInfoCollector.GetSystemSummary(userID)
			if err != nil {
				fmt.Println("error collecting system info:", err)
			} else {
				fmt.Println("systemInfo", sys)
				sender.SendSystemSummaryToAPI(sys)

			}
		}
	}
}
