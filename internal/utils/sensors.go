package utils

import (
	"runtime"

	"github.com/shirou/gopsutil/v3/host"
)

// IsSensorSupported returns true if the current OS supports temperature sensors
func IsSensorSupported() bool {
	// Temperature sensors generally not supported on Windows
	return runtime.GOOS != "windows"
}

// TemperatureInfo holds simplified sensor temperature data
type TemperatureInfo struct {
	SensorKey   string
	Temperature float64
	High        float64 // Optional: critical/high threshold
	Critical    float64 // Optional: critical temperature
}

// GetAllSensorTemperatures returns all available temperature readings
func GetAllSensorTemperatures() ([]TemperatureInfo, error) {
	if !IsSensorSupported() {
		return nil, nil
	}

	stats, err := host.SensorsTemperatures()
	if err != nil {
		return nil, err
	}

	var results []TemperatureInfo
	for _, stat := range stats {
		results = append(results, TemperatureInfo{
			SensorKey:   stat.SensorKey,
			Temperature: stat.Temperature,
			High:        stat.High,
			Critical:    stat.Critical,
		})
	}
	return results, nil
}

// GetSensorByName filters and returns sensor by sensor key name (e.g. "coretemp", "cpu-thermal")
func GetSensorByName(sensorName string) (*TemperatureInfo, error) {
	sensors, err := GetAllSensorTemperatures()
	if err != nil {
		return nil, err
	}

	for _, s := range sensors {
		if s.SensorKey == sensorName {
			return &s, nil
		}
	}
	return nil, nil // not found
}
