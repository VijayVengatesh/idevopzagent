package main

import (
	"fmt"
	"iDevopzAgent/internal/platform"
)

func main() {
	m := platform.GetMetricsCollector()

	fmt.Println("CPU  :", m.GetCPUUsage())
	fmt.Println("Memory:", m.GetMemoryUsage())
}
