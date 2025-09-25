// cmd/yourapp/main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"iDevopzAgent/configs"
	"iDevopzAgent/internal/healthreport"
	"iDevopzAgent/internal/metrics"
	"iDevopzAgent/internal/processdetails"
	"iDevopzAgent/internal/systeminfo"
	"iDevopzAgent/internal/utils"
	"iDevopzAgent/pkg/logger"
	"iDevopzAgent/security"
	"iDevopzAgent/sender"

	"github.com/gorilla/websocket"
)

type Process struct {
	PID    int    `json:"pid"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SocketMessage struct {
	MonitorID string  `json:"monitorId"`
	Action    string  `json:"action"`
	Process   Process `json:"process"`
}

type APIProcess struct {
	Result        string  `json:"result"`
	Table         int     `json:"table"`
	Time          string  `json:"_time"`
	CPUPercent    float64 `json:"cpu_percent"`
	HandleCount   int     `json:"handle_count"`
	Hostname      string  `json:"hostname"`
	MemoryPercent float64 `json:"memory_percent"`
	Name          string  `json:"name"`
	Path          string  `json:"path"`
	PID           int     `json:"pid"`
	Priority      string  `json:"priority"`
	ThreadCount   int     `json:"thread_count"`
	UserName      string  `json:"user_name"`
}

type ProcessUpdate struct {
	MonitorID string `json:"monitorId"`
	PID       int    `json:"pid"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Hostname  string `json:"hostname"`
}

var socketClient *websocket.Conn
var metricsEndpoint = configs.LoadConfig().APIEndpoint + "/metrics"
var processStatusMap = make(map[int]string)
var processNameMap = make(map[int]string)
var refreshProcessesChan = make(chan bool, 1)

func main() {
	// Initialize Socket.IO connection with retry
	go initSocketIOConnection()
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
	logger.Log.WithFields(map[string]interface{}{
		"userID":    userID,
		"machineID": machineID,
		"hostname":  hostname,
		"os":        os,
	}).Info("Startup details")

	startupPayload := map[string]string{
		"monitorId": userID,
		"hostname":  hostname,
		"machineId": machineID,
		"os":        os,
	}

	sender.SendStartupAPI(startupPayload)

	go collectMetrics(userID, machineID)

	// Fetch selected processes at startup
	go func() {
		time.Sleep(3 * time.Second) // Wait for socket connection
		processes, err := getSelectedProcesses(userID, machineID, hostname)
		if err == nil {
			logger.Log.WithField("count", len(processes)).Info("Initial selected processes loaded")
			for _, proc := range processes {
				processNameMap[proc.PID] = proc.Name
				if isProcessRunning(proc.PID, proc.Name) {
					processStatusMap[proc.PID] = "up"
				} else {
					processStatusMap[proc.PID] = "down"
				}
			}
		}
	}()

	go monitorProcesses(userID, machineID, hostname)
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
			if err != nil {
				logger.Log.WithError(err).Error("Failed to collect metrics")
			} else {
				logger.Log.WithField("metrics", y).Info("Metrics collected")

				sender.SendToMetricsAPI(y)
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
				logger.Log.WithError(err).Error("Error collecting health report")
			} else {
				logger.Log.WithField("healthReport", health).Info("Collected health report")
				sender.SendToHealthReportAPI(health)
			}
		}
	}
}

func collectProcessDetails(userID string, machineId string) {
	tickerAll := time.NewTicker(1 * time.Minute)
	tickerTop := time.NewTicker(10 * time.Second)

	defer tickerAll.Stop()
	defer tickerTop.Stop()

	processUtil := processdetails.GetProcessCollector()

	for {
		select {
		case <-tickerAll.C:
			if p, err := processUtil.ListAllProcesses(userID, machineId); err == nil {
				logger.Log.WithField("process_count", len(p)).Info("Collected process list, sending to API")
				sender.SendProcessList(p)
			} else {
				logger.Log.WithError(err).Error("Error collecting process list")
			}

		case <-tickerTop.C:
			if top5Cpu, err := processUtil.ListTop5CpuProcess(userID, machineId); err == nil {
				sender.Top5Cpu(top5Cpu)
			}

			if top5Mem, err := processUtil.ListTop5MemoryProcess(userID, machineId); err == nil {
				sender.Top5Memory(top5Mem)
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
				logger.Log.WithError(err).Error("Error collecting system info")
			} else {
				logger.Log.WithField("systemInfo", sys).Info("Collected system info")
				sender.SendSystemSummaryToAPI(sys)
			}
		}
	}
}

func initSocketIOConnection() {
	for {
		config := configs.LoadConfig()
		// Fix: Remove http/https prefix and use correct websocket URL format
		host := strings.TrimPrefix(strings.TrimPrefix(config.APIEndpoint, "http://"), "https://")
		u := url.URL{Scheme: "ws", Host: host, Path: "/socket.io/", RawQuery: "EIO=4&transport=websocket"}
		dialer := websocket.Dialer{HandshakeTimeout: 10 * time.Second}

		conn, _, err := dialer.Dial(u.String(), nil)
		if err != nil {
			logger.Log.WithError(err).Error("Failed to connect to Socket.IO, retrying in 5 seconds")
			time.Sleep(5 * time.Second)
			continue
		}

		socketClient = conn
		fmt.Println("📊 Connected to Socket.IO server")
		logger.Log.Info("Socket.IO connection established")

		// Fix: Correct namespace connection format
		namespaceMsg := "40/metrics"
		err = conn.WriteMessage(websocket.TextMessage, []byte(namespaceMsg))
		if err != nil {
			logger.Log.WithError(err).Error("Failed to join metrics namespace")
			conn.Close()
			socketClient = nil
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("Joined /metrics namespace")

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				logger.Log.WithError(err).Error("Socket.IO connection lost")
				conn.Close()
				socketClient = nil
				break
			}

			msgStr := string(message)
			if strings.Contains(msgStr, "selected-process-refresh") {
				logger.Log.Info("Received selected-process-refresh event, triggering refresh")
				select {
				case refreshProcessesChan <- true:
				default:
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func sendMetricsViaWebSocket(data interface{}) {
	if socketClient == nil {
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to marshal metrics data")
		return
	}

	message := fmt.Sprintf(`42/metrics,["metric-data",%s]`, string(jsonData))
	err = socketClient.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		logger.Log.WithError(err).Error("Failed to send metrics via Socket.IO")
		socketClient.Close()
		socketClient = nil
	}
}

func monitorProcesses(userID, machineID, hostname string) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-refreshProcessesChan:
			logger.Log.Info("Refetching selected processes due to add/remove event")
			time.Sleep(2 * time.Second)
			processes, err := getSelectedProcesses(userID, machineID, hostname)
			if err == nil {
				logger.Log.WithField("count", len(processes)).Info("Selected processes refetched")
				// Clear existing maps
				processStatusMap = make(map[int]string)
				processNameMap = make(map[int]string)
				// Initialize with actual running status
				for _, proc := range processes {
					processNameMap[proc.PID] = proc.Name
					currentStatus := "down"
					if isProcessRunning(proc.PID, proc.Name) {
						currentStatus = "up"
					}
					processStatusMap[proc.PID] = currentStatus
				}
			} else {
				logger.Log.WithError(err).Error("Failed to refetch selected processes")
			}
		case <-ticker.C:
			// Monitor status changes for tracked processes
			fmt.Printf("🔍 Checking %d selected processes...\n", len(processStatusMap))
			for pid, oldStatus := range processStatusMap {
				processName := processNameMap[pid]
				isRunning := isProcessRunning(pid, processName)
				currentStatus := "down"
				if isRunning {
					currentStatus = "up"
				}

				// Print status check for each process
				fmt.Printf("🔍 PID: %d, Name: %s, Running: %t, Status: %s -> %s\n",
					pid, processName, isRunning, oldStatus, currentStatus)

				// Only send update if status changed
				if currentStatus != oldStatus {
					processStatusMap[pid] = currentStatus
					fmt.Printf("🔄 Status changed for %s (PID: %d): %s -> %s\n",
						processName, pid, oldStatus, currentStatus)
					go sendProcessUpdate(userID, pid, processName, currentStatus, hostname)
					logger.Log.WithFields(map[string]interface{}{
						"pid":       pid,
						"name":      processName,
						"oldStatus": oldStatus,
						"newStatus": currentStatus,
					}).Info("Process status changed")
				}
			}
		}
	}
}

func getSelectedProcesses(userID, machineID, hostname string) ([]APIProcess, error) {
	config := configs.LoadConfig()
	apiURL := fmt.Sprintf("%s/api/vm/moniters/%s/%s/processes", config.APIEndpoint, userID, machineID)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Print API response
	fmt.Printf("📋 Selected Processes API Response: %s\n", string(body))
	logger.Log.WithField("response", string(body)).Info("Selected processes API response")

	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			Process []struct {
				PID    int    `json:"pid"`
				Name   string `json:"name"`
				Status string `json:"status"`
			} `json:"process"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("API request failed: %s", response.Message)
	}

	var processes []APIProcess
	for _, proc := range response.Data.Process {
		actualRunning := utils.IsProcessRunning(proc.PID, proc.Name)
		actualStatus := "down"
		if actualRunning {
			actualStatus = "up"
		}
		fmt.Printf("📊 Process: %s (PID: %d) | API Status: %s | Actual Status: %s | Running: %t\n",
			proc.Name, proc.PID, proc.Status, actualStatus, actualRunning)

		// Initialize maps with API status
		processNameMap[proc.PID] = proc.Name
		processStatusMap[proc.PID] = proc.Status

		// If API status doesn't match actual status, send update immediately
		if proc.Status != actualStatus {
			fmt.Printf("⚠️ Status mismatch detected! Sending update for %s (PID: %d): %s -> %s\n",
				proc.Name, proc.PID, proc.Status, actualStatus)
			go func(pid int, name, status string) {
				sendProcessUpdate(userID, pid, name, status, hostname)
				processStatusMap[pid] = status
			}(proc.PID, proc.Name, actualStatus)
		}

		processes = append(processes, APIProcess{
			PID:  proc.PID,
			Name: proc.Name,
		})
	}

	return processes, nil
}

func isProcessRunning(pid int, name string) bool {
	return utils.IsProcessRunning(pid, name)
}

func sendProcessUpdate(monitorID string, pid int, name, status, hostname string) {
	update := ProcessUpdate{
		MonitorID: monitorID,
		PID:       pid,
		Name:      name,
		Status:    status,
		Hostname:  hostname,
	}

	jsonData, err := json.Marshal(update)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to marshal process update")
		return
	}

	config := configs.LoadConfig()
	apiURL := fmt.Sprintf("%s/api/vm/moniters/%s/update/processes", config.APIEndpoint, monitorID)
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log.WithError(err).Error("Failed to create process update request")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to send process update")
		return
	}
	defer resp.Body.Close()

	logger.Log.WithFields(map[string]interface{}{
		"monitorId": monitorID,
		"pid":       pid,
		"name":      name,
		"status":    status,
		"hostname":  hostname,
	}).Info("Process update sent successfully")
}
