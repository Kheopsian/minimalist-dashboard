package services

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"minimalist-dashboard/internal/config"
	"minimalist-dashboard/internal/models"
	"minimalist-dashboard/internal/utils"
)

// StorageService gère les informations de stockage et réseau
type StorageService struct {
	config *config.Config
}

// NewStorageService crée une nouvelle instance du service de stockage
func NewStorageService(cfg *config.Config) *StorageService {
	return &StorageService{config: cfg}
}

// GetDiskInfo récupère les informations sur le disque
func (s *StorageService) GetDiskInfo() (models.DiskInfo, error) {
	monitorPath := s.config.GetMonitorPath()

	out, err := utils.RunCommand("df", monitorPath)
	if err != nil {
		return models.DiskInfo{}, err
	}

	lines := strings.Split(out, "\n")
	if len(lines) < 2 {
		return models.DiskInfo{}, fmt.Errorf("invalid df output")
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 6 {
		return models.DiskInfo{}, fmt.Errorf("invalid df fields")
	}

	percentStr := strings.TrimRight(fields[4], "%")
	percentNum, _ := strconv.ParseFloat(percentStr, 64)

	totalK, _ := strconv.ParseFloat(fields[1], 64)
	usedK, _ := strconv.ParseFloat(fields[2], 64)

	// Convertir de Kio à Tio (1024*1024*1024)
	kibToTb := 1024.0 * 1024.0 * 1024.0

	return models.DiskInfo{
		Total:      fmt.Sprintf("%.1f TB", totalK/kibToTb),
		Used:       fmt.Sprintf("%.1f TB", usedK/kibToTb),
		Free:       fields[3],
		Percent:    fields[4],
		PercentNum: percentNum,
		MountPoint: fields[5],
	}, nil
}

// GetNetCounters récupère les compteurs réseau bruts
func (s *StorageService) GetNetCounters(interfaceName string) (models.NetCounters, error) {
	content, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return models.NetCounters{}, err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines[2:] {
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		currentInterface := strings.TrimRight(fields[0], ":")

		if currentInterface == interfaceName {
			rx, _ := strconv.ParseFloat(fields[1], 64)
			tx, _ := strconv.ParseFloat(fields[9], 64)
			return models.NetCounters{RxBytes: rx, TxBytes: tx}, nil
		}
	}
	return models.NetCounters{}, fmt.Errorf("interface %s non trouvée dans /proc/net/dev", interfaceName)
}