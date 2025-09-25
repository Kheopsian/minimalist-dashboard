package services

import (
	"fmt"
	"time"

	"minimalist-dashboard/internal/config"
	"minimalist-dashboard/internal/models"
	"minimalist-dashboard/internal/utils"
)

// MetricsService orchestrates all services to collect metrics
type MetricsService struct {
	config          *config.Config
	systemService   *SystemService
	cpuService      *CPUService
	storageService  *StorageService
	zfsService      *ZFSService
	dockerService   *DockerService
	streamingService *StreamingService
}

// NewMetricsService creates a new metrics service instance
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

// CollectAllMetrics collects all system metrics
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

// CollectRealTimeMetrics collects real-time metrics (CPU and network)
func (m *MetricsService) CollectRealTimeMetrics(prevCPUTimes models.CPUTimes, prevNetCounters models.NetCounters, prevTime time.Time) (models.GlobalMetrics, models.CPUTimes, models.NetCounters, time.Time) {
	currentTime := time.Now()
	elapsedSeconds := currentTime.Sub(prevTime).Seconds()

	// CPU calculation
	currentCPUTimes, _ := m.cpuService.GetCPUTimes()
	deltaIdle := currentCPUTimes.Idle - prevCPUTimes.Idle
	deltaTotal := currentCPUTimes.Total - prevCPUTimes.Total
	cpuUsagePercent := 0.0
	if deltaTotal > 0 {
		cpuUsagePercent = (1.0 - deltaIdle/deltaTotal) * 100
	}
	tempStr, tempDeg := m.cpuService.GetCPUTemp()

	// Collect base metrics
	metrics := m.CollectAllMetrics()

	// Network calculation
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