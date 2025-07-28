// cmd/yourapp/main.go
package main

import (
	"encoding/json"
	"fmt"
	"iDevopzAgent/configs"
	"iDevopzAgent/internal/healthreport"
	"iDevopzAgent/internal/metrics"
	"iDevopzAgent/internal/processdetails"
	"iDevopzAgent/internal/systeminfo"
	"iDevopzAgent/internal/utilization"
	"iDevopzAgent/internal/utils"
)

func main() {

	userID, err := configs.LoadUserID()
	if err != nil || userID == "" {
		userID = configs.PromptAndSaveUserID()
	}
	// Metrics
	collector := metrics.GetCollector()
	y, err := collector.MetricsCollect(userID)
	fmt.Println("collected metrics", y, err)

	u := utilization.UtilizationCollector()

	cpuUtil, err := u.CpuUtilization(userID)
	if err != nil {
		fmt.Println("error collecting cpu utilization:", err)
	} else {
		fmt.Println("cpu utilization", cpuUtil)
	}

	memUtil, err := u.MemoryUtilization(userID)
	if err != nil {
		fmt.Println("error collecting memory utilization:", err)
	} else {
		fmt.Println("memory utilzation", memUtil)
	}

	diskUtil, err := u.DiskUtilization(userID)
	if err != nil {
		fmt.Println("error collecting disk utilization:", err)
	} else {
		fmt.Println("disk utilzation", diskUtil)
	}

	h := healthreport.GetHealthReportCollector()

	helthReportUtil, err := h.GenerateHealthReport(userID)

	if err != nil {
		fmt.Println("error collecting of healthReport", err)
	} else {
		fmt.Println("healthReport", helthReportUtil)
	}

	processUtil := processdetails.GetProcessCollector()

	p, err := processUtil.ListAllProcesses(userID)

	if err != nil {
		fmt.Println("error collecting of process details", err)
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling processes:", err)
	} else {
		fmt.Println("process details", string(data), len(p))
	}
	top5Cpu, err := processUtil.ListTop5CpuProcess(userID)
	if err != nil {
		fmt.Println("Error getting top 5 processes by CPU:", err)
	} else {
		data, _ := json.MarshalIndent(top5Cpu, "", "  ")
		fmt.Println("Top 5 processes by CPU:\n", string(data))
	}

	top5Mem, err := processUtil.ListTop5MemoryProcess(userID)
	if err != nil {
		fmt.Println("Error getting top 5 processes by Memory:", err)
	} else {
		data, _ := json.MarshalIndent(top5Mem, "", "  ")
		fmt.Println("Top 5 processes by Memory:\n", string(data))
	}

	count, err := utils.GetProcessCount()
	if err != nil {
		fmt.Println("Error get Process count Error ", err)
	} else {
		fmt.Println("Process count ", count)
	}

	//sytemInfo

	systemInfoCollector := systeminfo.GetSystemInfoCollector()

	sys, err := systemInfoCollector.GetSystemSummary(userID)

	if err != nil {
		fmt.Println("error collecting of systemInfo", err)
	} else {
		fmt.Println("systemInfo", sys)
	}

}
