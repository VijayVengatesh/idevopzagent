package sender

import (
	"fmt"
	"iDevopzAgent/configs"
	"iDevopzAgent/httpclient"
	"iDevopzAgent/models"
	"io"
)

func SendStartupAPI(payload map[string]string) {
	url := configs.LoadConfig().APIEndpoint + "/api/vm/moniters/create-update"

	resp, err := httpclient.SendPOST(url, payload)
	if err != nil {
		fmt.Println("Error sending startup API:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("Startup API sent successfully:", resp.Status)
	} else {
		fmt.Printf("Startup API failed.\nStatus: %s\nResponse: %s\n", resp.Status, string(body))
	}
}

func SendToMetricsAPI(metrics *models.Metrics) {

	url := configs.LoadConfig().APIEndpoint + "/api/go/system/metrics/create"

	fmt.Println(" Sending payload to:", url)

	resp, err := httpclient.SendPOST(url, metrics)
	if err != nil {
		fmt.Println(" Error sending data:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println(" Metrics sent! Status:", resp.Status)
	} else {
		fmt.Printf(" Failed to send metrics.\nStatus: %s\nResponse: %s\n", resp.Status, string(body))
	}
}

func SendToHealthReportAPI(healthReport *models.HealthReport) {

	url := configs.LoadConfig().APIEndpoint + "/api/go/system/health-report"

	resp, err := httpclient.SendPOST(url, healthReport)
	if err != nil {
		fmt.Println(" Error sending healthReport:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println(" healthReport sent! Status:", resp.Status)
	} else {
		fmt.Printf(" Failed to send healthReport.\nStatus: %s\nResponse: %s\n", resp.Status, string(body))
	}
}

func SendSystemSummaryToAPI(report *models.Systeminfo) {

	url := configs.LoadConfig().APIEndpoint + "/api/go/system/summary"

	resp, err := httpclient.SendPOST(url, report)
	if err != nil {
		fmt.Println(" Error sending system summary:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println(" System summary sent! Status:", resp.Status)
	} else {
		fmt.Printf(" Failed to send system summary.\nStatus: %s\nResponse: %s\n", resp.Status, string(body))
	}
}

// func SendLoadAverageToAPI(report *models.LoadAverageMetrics) {
// 	jsonData, _ := json.MarshalIndent(report, "", "  ")
// 	fmt.Println(" Sending LoadAverage payload:\n", string(jsonData))

// 	url := configs.LoadConfig().APIEndpoint + "/send-load-averge"

// 	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		fmt.Println(" Error sending load average:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)
// 	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
// 		fmt.Println(" Load average sent! Status:", resp.Status)
// 	} else {
// 		fmt.Printf(" Failed to send load average.\nStatus: %s\nResponse: %s\n", resp.Status, string(body))
// 	}
// }

func SendCpuUtilizationToAPI(report *models.CpuUtilization) {

	url := configs.LoadConfig().APIEndpoint + "/api/go/send-cpu-utilization"

	resp, err := httpclient.SendPOST(url, report)
	if err != nil {
		fmt.Println(" Error sending CPU utilization:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println(" CPU utilization sent! Status:", resp.Status)
	} else {
		fmt.Printf(" Failed to send CPU utilization.\nStatus: %s\nResponse: %s\n", resp.Status, string(body))
	}
}
func SendMemmoryUtilizationToAPI(report *models.MemoryUtilization) {

	url := configs.LoadConfig().APIEndpoint + "/api/go/send-memory-utilization"

	resp, err := httpclient.SendPOST(url, report)
	if err != nil {
		fmt.Println(" Error sending Memory utilization:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println(" Memory utilization sent! Status:", resp.Status)
	} else {
		fmt.Printf(" Failed to send Memory utilization.\nStatus: %s\nResponse: %s\n", resp.Status, string(body))
	}
}
func SendDiskUtilizationToAPI(report *models.DiskUtilization) {

	url := configs.LoadConfig().APIEndpoint + "/api/go/send-disk-utilization"

	resp, err := httpclient.SendPOST(url, report)
	if err != nil {
		fmt.Println(" Error sending disk utilization:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println(" Disk utilization sent! Status:", resp.Status)
	} else {
		fmt.Printf(" Failed to send disk utilization.\nStatus: %s\nResponse: %s\n", resp.Status, string(body))
	}
}
func SendProcessList(report []*models.ProcessInfo) {

	url := configs.LoadConfig().APIEndpoint + "/api/go/system/create-processes"

	resp, err := httpclient.SendPOST(url, report)
	fmt.Println("response", resp)
	if err != nil {
		fmt.Println(" Error sending process list :", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println(" Process List  Status:", resp.Status)
	} else {
		fmt.Printf(" Failed to Process info list .\nStatus: %s\nResponse: %s\n", resp.Status, string(body))
	}
}

func Top5Cpu(report []*models.Process) {

	url := configs.LoadConfig().APIEndpoint + "/api/go/system/processes/topcpu-create"

	resp, err := httpclient.SendPOST(url, report)
	if err != nil {
		fmt.Println(" Error sending Top 5 Cpu List :", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println(" Process Top 5 Cpu list:", resp.Status)
	} else {
		fmt.Printf("Process Top 5 Cpu list.\nStatus: %s\nResponse: %s\n", resp.Status, string(body))
	}
}

func Top5Memory(report []*models.Process) {

	url := configs.LoadConfig().APIEndpoint + "/api/go/system/processes/topmemory-create"

	resp, err := httpclient.SendPOST(url, report)
	if err != nil {
		fmt.Println(" Error sending Top 5 Memory List :", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println(" Process Top 5 Memory list:", resp.Status)
	} else {
		fmt.Printf("Process Top 5 Memory list.\nStatus: %s\nResponse: %s\n", resp.Status, string(body))
	}
}
