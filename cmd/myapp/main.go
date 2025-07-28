// cmd/yourapp/main.go
package main

import (
	"fmt"
	"iDevopzAgent/configs"
	"iDevopzAgent/internal/metrics"
	"iDevopzAgent/internal/utilization"
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

}
