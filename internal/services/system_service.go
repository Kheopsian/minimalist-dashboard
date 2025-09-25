package services

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"minimalist-dashboard/internal/models"
	"minimalist-dashboard/internal/utils"
)

// SystemService gère les informations système
type SystemService struct{}

// NewSystemService crée une nouvelle instance du service système
func NewSystemService() *SystemService {
	return &SystemService{}
}

// GetSystemInfo récupère les informations système générales
func (s *SystemService) GetSystemInfo() models.SystemInfo {
	uptimeContent, err := os.ReadFile("/proc/uptime")
	var uptimeSeconds float64
	if err == nil {
		uptimeSeconds, _ = strconv.ParseFloat(strings.Fields(string(uptimeContent))[0], 64)
	} else {
		log.Printf("Erreur lecture uptime: %v", err)
	}
	days := int(uptimeSeconds) / (60 * 60 * 24)
	hours := (int(uptimeSeconds) / (60 * 60)) % 24
	minutes := (int(uptimeSeconds) / 60) % 60
	uptimeFormatted := fmt.Sprintf("%dd %dh %dm", days, hours, minutes)

	osContent, err := os.ReadFile("/etc/os-release")
	var osName string
	if err == nil {
		for _, line := range strings.Split(string(osContent), "\n") {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				osName = strings.Trim(strings.Split(line, "=")[1], `"`)
			}
		}
	} else {
		log.Printf("Erreur lecture os-release: %v", err)
	}
	if osName == "" {
		osName = "Unraid OS"
	}

	kernelOut, _ := utils.RunCommand("uname", "-r")

	cpuContent, err := os.ReadFile("/proc/cpuinfo")
	var cpuModel string
	if err == nil {
		for _, line := range strings.Split(string(cpuContent), "\n") {
			if strings.HasPrefix(line, "model name") {
				cpuModel = strings.TrimSpace(strings.Split(line, ":")[1])
				break
			}
		}
	} else {
		log.Printf("Erreur lecture cpuinfo: %v", err)
	}

	return models.SystemInfo{
		OS:     osName,
		Kernel: kernelOut,
		CPU:    cpuModel,
		Uptime: uptimeFormatted,
	}
}