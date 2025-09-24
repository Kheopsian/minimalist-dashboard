
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
	out, err := runCommand("cat", "/proc/meminfo")
	if err != nil {
		log.Printf("Erreur getRAMInfo: %v", err)
		return RAMInfo{}
	}

	var memTotal, memAvailable float64
	lines := strings.Split(out, "\n")

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
		Used:    fmt.Sprintf("%.1f GB", used/1024/1024),
		Total:   fmt.Sprintf("%.1f GB", memTotal/1024/1024),
		Percent: fmt.Sprintf("%.1f%%", percent),
		PercentNum: percent,
	}
}


func getSystemInfo() SystemInfo {
	// Uptime : /proc/uptime donne l'uptime en secondes (le premier nombre)
	uptimeOut, _ := runCommand("cat", "/proc/uptime")
	var uptimeSeconds float64
	if uptimeOut != "" {
		uptimeSeconds, _ = strconv.ParseFloat(strings.Fields(uptimeOut)[0], 64)
	}

	// Formatage de l'uptime en jours, heures, minutes
	days := int(uptimeSeconds) / (60 * 60 * 24)
	hours := (int(uptimeSeconds) / (60 * 60)) % 24
	minutes := (int(uptimeSeconds) / 60) % 60
	uptimeFormatted := fmt.Sprintf("%dd %dh %dm", days, hours, minutes)

	// OS : On tente toujours /etc/os-release, c'est le standard moderne
	osOut, _ := runCommand("cat", "/etc/os-release")
	var osName string
	for _, line := range strings.Split(osOut, "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			osName = strings.Trim(strings.Split(line, "=")[1], `"`)
		}
	}
    if osName == "" {
        osName = "Unraid OS"
    }

	kernelOut, _ := runCommand("uname", "-r")

	cpuOut, _ := runCommand("cat", "/proc/cpuinfo")
	var cpuModel string
	for _, line := range strings.Split(cpuOut, "\n") {
		if strings.HasPrefix(line, "model name") {
			cpuModel = strings.TrimSpace(strings.Split(line, ":")[1])
			break
		}
	}

	return SystemInfo{
		OS:     osName,
		Kernel: kernelOut,
		CPU:    cpuModel,
		Uptime: uptimeFormatted,
	}
}

// collectAllMetrics gathers all system metrics.
func collectAllMetrics() GlobalMetrics {
	metrics := GlobalMetrics{
		RAM: getRAMInfo(),
		System: getSystemInfo(),
		Streaming: getStreamingInfo(),
		
		// Mocked data for now
		CPU: CPUInfo{Usage: "72%", Temp: "51.2°C", TempDeg: 51.2},
		Disk: DiskInfo{Total: "8 TB", Used: "4.5 TB", Free: "3.5 TB", Percent: "56.2%", PercentNum: 56.2},
		Net: NetTraffic{In: "5.7 MB/s", Out: "1.2 MB/s"},
		ZFSConfig: ZFSConfig{PoolName: "rpool", PoolStatus: "ONLINE", DataVdevs: []ZPoolVdev{{Name: "raidz1-0", Status: "ONLINE", Devices: []string{"disk-01", "disk-02", "disk-03", "disk-04"}},{Name: "raidz1-1", Status: "ONLINE", Devices: []string{"disk-05", "disk-06", "disk-07", "disk-08"}}}, CacheVdev: &ZPoolVdev{Name: "L2ARC", Status: "ONLINE", Devices: []string{"nvme-Samsung-970-Evo"}}},
		ARCCache: ARCCache{ARCSize: "66.1 GB", ARCTargetSize: "65.0 GB", ARCHitRate: "98.1%", ARCHitRateNum: 98.1, L2ARCSize: "485.0 GB", L2ARCHitRate: "70.5 %"},
		Docker: DockerInfo{Containers: 25, Images: 132, Volumes: 38},
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