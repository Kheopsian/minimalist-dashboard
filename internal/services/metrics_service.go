package services

import (
	"fmt"
	"time"

	"minimalist-dashboard/internal/config"
	"minimalist-dashboard/internal/models"
	"minimalist-dashboard/internal/utils"
)

// MetricsService orchestre tous les services pour collecter les métriques
type MetricsService struct {
	config          *config.Config
	systemService   *SystemService
	cpuService      *CPUService
	storageService  *StorageService
	zfsService      *ZFSService
	dockerService   *DockerService
	streamingService *StreamingService
}

// NewMetricsService crée une nouvelle instance du service de métriques
func NewMetricsService(cfg *config.Config) *MetricsService {
	return &MetricsService{
		config:           cfg,
		systemService:    NewSystemService(),
		cpuService:       NewCPUService(),
		storageService:   NewStorageService(cfg),
		zfsService:       NewZFSService(),
		dockerService:    NewDockerService(),
		streamingService: NewStreamingService(cfg),
	}
}

// CollectAllMetrics collecte toutes les métriques système
func (m *MetricsService) CollectAllMetrics() models.GlobalMetrics {
	diskInfo, _ := m.storageService.GetDiskInfo()

	metrics := models.GlobalMetrics{
		RAM:       m.cpuService.GetRAMInfo(),
		System:    m.systemService.GetSystemInfo(),
		Streaming: m.streamingService.GetStreamingInfo(),
		Docker:    m.dockerService.GetDockerInfo(),
		ZFSConfig: m.zfsService.GetZFSConfig(),
		ARCCache:  m.zfsService.GetARCCacheInfo(),
		Disk:      diskInfo,
	}
	return metrics
}

// CollectRealTimeMetrics collecte les métriques en temps réel (CPU et réseau)
func (m *MetricsService) CollectRealTimeMetrics(prevCPUTimes models.CPUTimes, prevNetCounters models.NetCounters, prevTime time.Time) (models.GlobalMetrics, models.CPUTimes, models.NetCounters, time.Time) {
	currentTime := time.Now()
	elapsedSeconds := currentTime.Sub(prevTime).Seconds()

	// Calcul CPU
	currentCPUTimes, _ := m.cpuService.GetCPUTimes()
	deltaIdle := currentCPUTimes.Idle - prevCPUTimes.Idle
	deltaTotal := currentCPUTimes.Total - prevCPUTimes.Total
	cpuUsagePercent := 0.0
	if deltaTotal > 0 {
		cpuUsagePercent = (1.0 - deltaIdle/deltaTotal) * 100
	}
	tempStr, tempDeg := m.cpuService.GetCPUTemp()

	// Collecte des métriques de base
	metrics := m.CollectAllMetrics()

	// Calcul réseau
	currentNetCounters, err := m.storageService.GetNetCounters(m.config.NetInterface)
	if err == nil {
		deltaRx := currentNetCounters.RxBytes - prevNetCounters.RxBytes
		deltaTx := currentNetCounters.TxBytes - prevNetCounters.TxBytes

		rxSpeed := deltaRx / elapsedSeconds
		txSpeed := deltaTx / elapsedSeconds

		metrics.Net = models.NetTraffic{
			In:  utils.FormatSpeed(rxSpeed),
			Out: utils.FormatSpeed(txSpeed),
		}
		prevNetCounters = currentNetCounters
	}

	metrics.CPU = models.CPUInfo{
		Usage:   fmt.Sprintf("%.0f%%", cpuUsagePercent),
		Temp:    tempStr,
		TempDeg: tempDeg,
	}

	return metrics, currentCPUTimes, prevNetCounters, currentTime
}