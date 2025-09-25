package services

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"minimalist-dashboard/internal/models"
)

// ZFSService manages ZFS information
type ZFSService struct{}

// NewZFSService creates a new ZFS service instance
func NewZFSService() *ZFSService {
	return &ZFSService{}
}

// GetZFSConfig retrieves ZFS configuration
func (z *ZFSService) GetZFSConfig() models.ZFSConfig {
	content, err := os.ReadFile("/app/zpool_status.txt")
	if err != nil {
		log.Printf("Error getZFSConfig: unable to read /app/zpool_status.txt: %v", err)
		return models.ZFSConfig{}
	}

	out := string(content)

	config := models.ZFSConfig{}
	var dataVdevs []models.ZPoolVdev
	var cacheVdev *models.ZPoolVdev

	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "pool:":
			config.PoolName = fields[1]
		case "state:":
			config.PoolStatus = fields[1]
		}
	}

	inConfigSection := false
	var lastVdev *models.ZPoolVdev

	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "NAME") && strings.Contains(line, "STATE") {
			inConfigSection = true
			continue
		}
		if strings.HasPrefix(line, "errors:") {
			break
		}
		if !inConfigSection || len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), config.PoolName) {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue // Ignore completely empty lines
		}

		deviceName := fields[0]
		deviceStatus := "" // Default status
		if len(fields) > 1 {
			deviceStatus = fields[1]
		}

		if strings.Contains(deviceName, "raidz") || strings.Contains(deviceName, "mirror") {
			dataVdevs = append(dataVdevs, models.ZPoolVdev{Name: deviceName, Status: deviceStatus})
			lastVdev = &dataVdevs[len(dataVdevs)-1]
		} else if deviceName == "cache" {
			cacheVdev = &models.ZPoolVdev{Name: deviceName, Status: deviceStatus}
			lastVdev = cacheVdev
		} else if lastVdev != nil {
			lastVdev.Devices = append(lastVdev.Devices, deviceName)
		}
	}

	config.DataVdevs = dataVdevs
	config.CacheVdev = cacheVdev

	return config
}

// GetARCCacheInfo retrieves ARC cache information
func (z *ZFSService) GetARCCacheInfo() models.ARCCache {
	content, err := os.ReadFile("/proc/spl/kstat/zfs/arcstats")
	if err != nil {
		log.Printf("Error getARCCacheInfo: unable to read /proc/spl/kstat/zfs/arcstats: %v", err)
		return models.ARCCache{}
	}

	stats := make(map[string]float64)
	lines := strings.Split(string(content), "\n")

	// The 3rd line contains headers, data starts after
	if len(lines) < 3 {
		return models.ARCCache{}
	}

	for _, line := range lines[2:] {
		fields := strings.Fields(line)
		if len(fields) == 3 {
			// Format is: stat_name type value
			key := fields[0]
			value, _ := strconv.ParseFloat(fields[2], 64)
			stats[key] = value
		}
	}

	arcHitrate := 0.0
	if (stats["hits"] + stats["misses"]) > 0 {
		arcHitrate = (stats["hits"] / (stats["hits"] + stats["misses"])) * 100
	}

	l2arcHitrate := 0.0
	if (stats["l2_hits"] + stats["l2_misses"]) > 0 {
		l2arcHitrate = (stats["l2_hits"] / (stats["l2_hits"] + stats["l2_misses"])) * 100
	}

	return models.ARCCache{
		ARCSize:       fmt.Sprintf("%.1f GB", stats["size"]/1024/1024/1024),
		ARCMaxSize:    fmt.Sprintf("%.1f GB", stats["c_max"]/1024/1024/1024),
		ARCHitRate:    fmt.Sprintf("%.1f%%", arcHitrate),
		ARCHitRateNum: arcHitrate,
		L2ARCSize:     fmt.Sprintf("%.1f GB", stats["l2_size"]/1024/1024/1024),
		L2ARCHitRate:  fmt.Sprintf("%.1f%%", l2arcHitrate),
	}
}