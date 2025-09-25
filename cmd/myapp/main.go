// cmd/yourapp/main.go
package main

import (
	"fmt"
	"time"

	"iDevopzAgent/configs"
	"iDevopzAgent/internal/healthreport"
	"iDevopzAgent/internal/metrics"
	"iDevopzAgent/internal/processdetails"
	"iDevopzAgent/internal/systeminfo"
	"iDevopzAgent/internal/utilization"
	"iDevopzAgent/internal/utils"
	"iDevopzAgent/security"
	"iDevopzAgent/sender"
)

func main() {
	userID, machineID, err := configs.LoadUserID()
	if err != nil || userID == "" {
		encUserID, encMachineID := configs.PromptAndSaveUserID()

		// Decrypt immediately for runtime use
		decUserID, _ := security.Decrypt(encUserID)
		decMachineID, _ := security.Decrypt(encMachineID)

		userID = decUserID
		machineID = decMachineID
	}

	hostname, _ := utils.GetHostName()
	os := utils.GetOS()

	fmt.Println("UserID:", userID)
	fmt.Println("MachineID:", machineID)
	fmt.Println("Hostname:", hostname)
	fmt.Println("OS:", os)

	// ----------------------------
	// Send startup API once
	// ----------------------------
	startupPayload := map[string]string{
		"monitorId": userID,
		"hostname":  hostname,
		"machineId": machineID,
		"os":        os,
	}

	// Call startup API (errors are handled inside the function)

	sender.SendStartupAPI(startupPayload)

	go collectMetrics(userID, machineID)
	// go collectUtilization(userID, machineID)
	go collectHealthReport(userID, machineID)
	go collectProcessDetails(userID, machineID)
	go collectSystemInfo(userID, machineID)

	// Prevent the main function from exiting
	select {}
}

func collectMetrics(userID string, machineId string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	collector := metrics.GetCollector()

	for {
		select {
		case <-ticker.C:
			y, err := collector.MetricsCollect(userID, machineId)
			sender.SendToMetricsAPI(y)
			fmt.Println("collected metrics", y, err)
		}
	}
}

func collectUtilization(userID string, machineId string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	u := utilization.UtilizationCollector()

	for {
		select {
		case <-ticker.C:
			if cpuUtil, err := u.CpuUtilization(userID, machineId); err == nil {
				fmt.Println("cpu utilization", cpuUtil)
				sender.SendCpuUtilizationToAPI(cpuUtil)

			}
			if memUtil, err := u.MemoryUtilization(userID, machineId); err == nil {
				fmt.Println("memory utilization", memUtil)
				sender.SendMemmoryUtilizationToAPI(memUtil)

			}
			if diskUtil, err := u.DiskUtilization(userID, machineId); err == nil {
				fmt.Println("disk utilization", diskUtil)
				sender.SendDiskUtilizationToAPI(diskUtil)

			}
		}
	}
}

func collectHealthReport(userID string, machineId string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	h := healthreport.GetHealthReportCollector()

	for {
		select {
		case <-ticker.C:
			health, err := h.GenerateHealthReport(userID, machineId)
			if err != nil {
				fmt.Println("Error collecting healthReport:", err)
			} else {
				fmt.Println("healthReport", health)
				sender.SendToHealthReportAPI(health)

			}
		}
	}
}

func collectProcessDetails(userID string, machineId string) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	processUtil := processdetails.GetProcessCollector()

	for {
		select {
		case <-ticker.C:
			if p, err := processUtil.ListAllProcesses(userID, machineId); err == nil {
				fmt.Printf("Collected %d processes, sending to API\n", len(p))
				sender.SendProcessList(p)
			} else {
				fmt.Println("Error collecting process list:", err)
			}

			if top5Cpu, err := processUtil.ListTop5CpuProcess(userID, machineId); err == nil {
				fmt.Printf("Collected top 5 CPU processes, sending to API\n")
				sender.Top5Cpu(top5Cpu)
			} else {
				fmt.Println("Error collecting top 5 CPU processes:", err)
			}

			if top5Mem, err := processUtil.ListTop5MemoryProcess(userID, machineId); err == nil {
				fmt.Printf("Collected top 5 memory processes, sending to API\n")
				sender.Top5Memory(top5Mem)
			} else {
				fmt.Println("Error collecting top 5 memory processes:", err)
			}

			if count, err := utils.GetProcessCount(); err == nil {
				fmt.Println("Process count:", count)
			}
		}
	}
}

func collectSystemInfo(userID string, machineId string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	systemInfoCollector := systeminfo.GetSystemInfoCollector()

	for {
		select {
		case <-ticker.C:
			sys, err := systemInfoCollector.GetSystemSummary(userID, machineId)
			if err != nil {
				fmt.Println("error collecting system info:", err)
			} else {
				fmt.Println("systemInfo", sys)
				sender.SendSystemSummaryToAPI(sys)

			}
		}
	}
}
