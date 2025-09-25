package services

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"minimalist-dashboard/internal/models"
)

// CPUService manages CPU and memory information
type CPUService struct{}

// NewCPUService creates a new CPU service instance
func NewCPUService() *CPUService {
	return &CPUService{}
}

// GetCPUTimes retrieves current CPU times for usage calculation
func (c *CPUService) GetCPUTimes() (models.CPUTimes, error) {
	content, err := os.ReadFile("/proc/stat")
	if err != nil {
		return models.CPUTimes{}, err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") { // The first "cpu" line is the aggregate
			fields := strings.Fields(line)[1:]
			var total float64
			var idle float64
			for i, field := range fields {
				val, _ := strconv.ParseFloat(field, 64)
				total += val
				if i == 3 || i == 4 { // The idle and iowait fields
					idle += val
				}
			}
			return models.CPUTimes{Idle: idle, Total: total}, nil
		}
	}
	return models.CPUTimes{}, fmt.Errorf("'cpu' line not found in /proc/stat")
}

// GetCPUTemp retrieves CPU temperature
func (c *CPUService) GetCPUTemp() (string, float64) {
	// CPU temperature is often in this file, in millidegrees Celsius
	content, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return "N/A", 0
	}
	temp, _ := strconv.ParseFloat(strings.TrimSpace(string(content)), 64)
	temp /= 1000 // Convert to degrees
	return fmt.Sprintf("%.1f°C", temp), temp
}

// GetRAMInfo retrieves RAM memory information
func (c *CPUService) GetRAMInfo() models.RAMInfo {
	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return models.RAMInfo{}
	}

	var memTotal, memAvailable float64
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			memTotal, _ = strconv.ParseFloat(fields[1], 64)
		case "MemAvailable:":
			memAvailable, _ = strconv.ParseFloat(fields[1], 64)
		}
	}

	used := memTotal - memAvailable
	percent := 0.0
	if memTotal > 0 {
		percent = (used / memTotal) * 100
	}

	return models.RAMInfo{
		Used:       fmt.Sprintf("%.1f GB", used/1024/1024),
		Total:      fmt.Sprintf("%.1f GB", memTotal/1024/1024),
		Percent:    fmt.Sprintf("%.1f%%", percent),
		PercentNum: percent,
	}
}