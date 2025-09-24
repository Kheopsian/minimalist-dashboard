
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
	"encoding/json"
	"os/exec"
	"strings"
	"strconv"
	"os"
)


type CPUInfo struct {
	Usage    string `json:"usage"`
	Temp     string `json:"temp"`
	TempDeg  float64 `json:"tempDeg"`
}

type RAMInfo struct {
	Used     string `json:"used"`
	Total    string `json:"total"`
	Percent  string `json:"percent"`
	PercentNum float64 `json:"percentNum"`
}

type DiskInfo struct {
	Total    string `json:"total"`
	Used     string `json:"used"`
	Free     string `json:"free"`
	Percent  string `json:"percent"`
	PercentNum float64 `json:"percentNum"`
}

type NetTraffic struct {
	In  string `json:"in"`
	Out string `json:"out"`
}

type ZPoolVdev struct {
	Name    string   `json:"name"`
	Status  string   `json:"status"`
	Devices []string `json:"devices"`
}

type ZFSConfig struct {
	PoolName    string      `json:"poolName"`
	PoolStatus  string      `json:"poolStatus"`
	DataVdevs   []ZPoolVdev `json:"dataVdevs"`
	CacheVdev   *ZPoolVdev  `json:"cacheVdev,omitempty"` // Pointer to handle absence of L2ARC
}

type ARCCache struct {
	ARCSize       string  `json:"arcSize"`
	ARCTargetSize string  `json:"arcTargetSize"`
	ARCHitRate    string  `json:"arcHitRate"`
	ARCHitRateNum float64 `json:"arcHitRateNum"`
	L2ARCSize     string  `json:"l2arcSize"`
	L2ARCHitRate  string  `json:"l2arcHitRate"`
}

type SystemInfo struct {
	OS     string `json:"os"`
	Kernel string `json:"kernel"`
	CPU    string `json:"cpu"`
	Uptime string `json:"uptime"`
}

type DockerInfo struct {
	Containers int `json:"containers"`
	Images     int `json:"images"`
	Volumes    int `json:"volumes"`
}

type StreamingInfo struct {
	Films      int `json:"films"`
	Series     int `json:"series"`
	Animes     int `json:"animes"`
	Playing    int `json:"playing"`
	Transcoding int `json:"transcoding"`
}

// GlobalMetrics groups all the metrics to be sent to the frontend.
type GlobalMetrics struct {
	CPU       CPUInfo       `json:"cpu"`
	RAM       RAMInfo       `json:"ram"`
	Disk      DiskInfo      `json:"disk"`
	Net       NetTraffic    `json:"net"`
	ZFSConfig ZFSConfig     `json:"zfsConfig"`
	ARCCache  ARCCache      `json:"arcCache"`
	System    SystemInfo    `json:"system"`
	Docker    DockerInfo    `json:"docker"`
	Streaming StreamingInfo `json:"streaming"`
}

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing '%s %v': %w", name, args, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func getStreamingInfo() StreamingInfo {
	countItemsInDir := func(path string) int {
		if path == "" {
			return 0
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			log.Printf("Error reading directory %s: %v", path, err)
			return 0
		}
		return len(entries)
	}

	filmsPath := os.Getenv("PATH_FILMS")
	seriesPath := os.Getenv("PATH_SERIES")
	animesPath := os.Getenv("PATH_ANIMES")

	// Playback/transcoding stats are mocked for now
	return StreamingInfo{
		Films:      countItemsInDir(filmsPath),
		Series:     countItemsInDir(seriesPath),
		Animes:     countItemsInDir(animesPath),
		Playing:    2,  // TODO: Replace with real data from Plex/Jellyfin API
		Transcoding: 1,  // TODO: Replace with real data from Plex/Jellyfin API
	}
}

func getRAMInfo() RAMInfo {
	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		log.Printf("Erreur getRAMInfo: impossible de lire /proc/meminfo: %v", err)
		return RAMInfo{}
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

	return RAMInfo{
		Used:       fmt.Sprintf("%.1f GB", used/1024/1024),
		Total:      fmt.Sprintf("%.1f GB", memTotal/1024/1024),
		Percent:    fmt.Sprintf("%.1f%%", percent),
		PercentNum: percent,
	}
}


func getSystemInfo() SystemInfo {
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

	kernelOut, _ := runCommand("uname", "-r")

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

	return SystemInfo{
		OS:     osName,
		Kernel: kernelOut,
		CPU:    cpuModel,
		Uptime: uptimeFormatted,
	}
}

func getARCCacheInfo() ARCCache {
	content, err := os.ReadFile("/proc/spl/kstat/zfs/arcstats")
	if err != nil {
		log.Printf("Erreur getARCCacheInfo: impossible de lire /proc/spl/kstat/zfs/arcstats: %v", err)
		return ARCCache{}
	}

	stats := make(map[string]float64)
	lines := strings.Split(string(content), "\n")

	// La 3ème ligne contient les en-têtes, les données commencent après
	if len(lines) < 3 {
		return ARCCache{}
	}

	for _, line := range lines[2:] {
		fields := strings.Fields(line)
		if len(fields) == 3 {
			// Le format est : nom_de_la_stat type valeur
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

	return ARCCache{
		ARCSize:       fmt.Sprintf("%.1f GB", stats["size"]/1024/1024/1024),
		ARCTargetSize: fmt.Sprintf("%.1f GB", stats["c"]/1024/1024/1024),
		ARCHitRate:    fmt.Sprintf("%.1f%%", arcHitrate),
		ARCHitRateNum: arcHitrate,
		L2ARCSize:     fmt.Sprintf("%.1f GB", stats["l2_size"]/1024/1024/1024),
		L2ARCHitRate:  fmt.Sprintf("%.1f%%", l2arcHitrate),
	}
}

func getZFSConfig() ZFSConfig {
	log.Println("--- Début getZFSConfig (Debug v2) ---")

	content, err := os.ReadFile("/app/zpool_status.txt")
	if err != nil {
		log.Printf("[ERREUR] Impossible de lire /app/zpool_status.txt: %v", err)
		return ZFSConfig{}
	}

	out := string(content)
	config := ZFSConfig{}
	var dataVdevs []ZPoolVdev
	var cacheVdev *ZPoolVdev
	
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 { continue }
		switch fields[0] {
		case "pool:":
			config.PoolName = fields[1]
		case "state:":
			config.PoolStatus = fields[1]
		}
	}

	inConfigSection := false
	var lastVdev *ZPoolVdev

	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "NAME") && strings.Contains(line, "STATE") {
			inConfigSection = true
			continue
		}
		if strings.HasPrefix(line, "errors:") { break }
		if !inConfigSection || len(strings.TrimSpace(line)) == 0 { continue }
		if strings.HasPrefix(strings.TrimSpace(line), config.PoolName) { continue }
		
		log.Printf("Ligne en cours de traitement: '%s'", strings.TrimSpace(line))

		fields := strings.Fields(line)
		if len(fields) < 2 { continue }
		
		deviceName := fields[0]
		deviceStatus := fields[1]

		if strings.Contains(deviceName, "raidz") || strings.Contains(deviceName, "mirror") {
			log.Printf("  -> Détecté comme VDEV de données: %s", deviceName)
			dataVdevs = append(dataVdevs, ZPoolVdev{Name: deviceName, Status: deviceStatus})
			lastVdev = &dataVdevs[len(dataVdevs)-1]
			log.Printf("  ==> 'Parent' actuel défini sur: %s", lastVdev.Name)
		} else if deviceName == "cache" {
			log.Printf("  -> Détecté comme VDEV de cache: %s", deviceName)
			cacheVdev = &ZPoolVdev{Name: deviceName, Status: deviceStatus}
			lastVdev = cacheVdev
			log.Printf("  ==> 'Parent' actuel défini sur: %s", lastVdev.Name)
		} else if lastVdev != nil {
			log.Printf("  -> Détecté comme disque: %s. Ajout au parent '%s'", deviceName, lastVdev.Name)
			lastVdev.Devices = append(lastVdev.Devices, deviceName)
		} else {
            log.Printf("  -> Ligne ignorée (pas de parent trouvé): %s", deviceName)
        }
	}
	
	config.DataVdevs = dataVdevs
	config.CacheVdev = cacheVdev
	
	log.Println("--- Fin getZFSConfig (Debug v2) ---")
	return config
}

func getDockerInfo() DockerInfo {
	containersOut, _ := runCommand("docker", "ps", "--format", "{{.ID}}")
	imagesOut, _ := runCommand("docker", "images", "--format", "{{.ID}}")
	volumesOut, _ := runCommand("docker", "volume", "ls", "--format", "{{.Name}}")

	countLines := func(output string) int {
		if output == "" { return 0 }
		return len(strings.Split(output, "\n"))
	}
	
	return DockerInfo{
		Containers: countLines(containersOut),
		Images:     countLines(imagesOut),
		Volumes:    countLines(volumesOut),
	}
}

// collectAllMetrics gathers all system metrics.
func collectAllMetrics() GlobalMetrics {
	metrics := GlobalMetrics{
		RAM: getRAMInfo(),
		System: getSystemInfo(),
		Streaming: getStreamingInfo(),
		Docker: getDockerInfo(),
		ZFSConfig: getZFSConfig(),
		ARCCache:  getARCCacheInfo(), 
		
		// Mocked data for now
		CPU: CPUInfo{Usage: "72%", Temp: "51.2°C", TempDeg: 51.2},
		Disk: DiskInfo{Total: "8 TB", Used: "4.5 TB", Free: "3.5 TB", Percent: "56.2%", PercentNum: 56.2},
		Net: NetTraffic{In: "5.7 MB/s", Out: "1.2 MB/s"},
	}
	return metrics
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	fmt.Println("New client connected to WebSocket")

	for {
		metrics := collectAllMetrics()

		jsonMessage, err := json.Marshal(metrics)
		if err != nil {
			log.Printf("JSON encoding error: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = ws.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			log.Printf("Send error: %v", err)
			break
		}

		time.Sleep(2 * time.Second)
	}
}

func main() {
	port := 8080

	fileServer := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fileServer)
	http.HandleFunc("/ws", handleConnections)

	fmt.Printf("Server started. Go to http://localhost:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}