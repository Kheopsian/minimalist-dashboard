package services

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"minimalist-dashboard/internal/models"
)

// CPUService gère les informations CPU et mémoire
type CPUService struct{}

// NewCPUService crée une nouvelle instance du service CPU
func NewCPUService() *CPUService {
	return &CPUService{}
}

// GetCPUTimes récupère les temps CPU actuels pour le calcul d'utilisation
func (c *CPUService) GetCPUTimes() (models.CPUTimes, error) {
	content, err := os.ReadFile("/proc/stat")
	if err != nil {
		return models.CPUTimes{}, err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") { // La première ligne "cpu" est l'agrégat
			fields := strings.Fields(line)[1:]
			var total float64
			var idle float64
			for i, field := range fields {
				val, _ := strconv.ParseFloat(field, 64)
				total += val
				if i == 3 || i == 4 { // Les champs idle et iowait
					idle += val
				}
			}
			return models.CPUTimes{Idle: idle, Total: total}, nil
		}
	}
	return models.CPUTimes{}, fmt.Errorf("ligne 'cpu' non trouvée dans /proc/stat")
}

// GetCPUTemp récupère la température du CPU
func (c *CPUService) GetCPUTemp() (string, float64) {
	// La température CPU est souvent dans ce fichier, en millidegrés Celsius
	content, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return "N/A", 0
	}
	temp, _ := strconv.ParseFloat(strings.TrimSpace(string(content)), 64)
	temp /= 1000 // Convertir en degrés
	return fmt.Sprintf("%.1f°C", temp), temp
}

// GetRAMInfo récupère les informations sur la mémoire RAM
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