// main.go

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

// Définition des structs pour organiser nos métriques
// Les balises `json:"..."` indiquent le nom du champ en JSON.
// `omitempty` signifie que si le champ est vide, il n'apparaîtra pas dans le JSON.

type CPUInfo struct {
	Usage    string `json:"usage"` // Ex: "65%"
	Temp     string `json:"temp"`  // Ex: "48.5°C"
	TempDeg  float64 `json:"tempDeg"` // Pour la jauge, la valeur numérique
}

type RAMInfo struct {
	Used     string `json:"used"`     // Ex: "35.8 GB"
	Total    string `json:"total"`    // Ex: "128 GB"
	Percent  string `json:"percent"`  // Ex: "28%"
	PercentNum float64 `json:"percentNum"` // Pour la barre de progression, la valeur numérique
}

type DiskInfo struct {
	Total    string `json:"total"`    // Ex: "8 TB"
	Used     string `json:"used"`     // Ex: "4.2 TB"
	Free     string `json:"free"`     // Ex: "3.8 TB"
	Percent  string `json:"percent"`  // Ex: "52.5%"
	PercentNum float64 `json:"percentNum"` // Pour la barre de progression, la valeur numérique
}

type NetTraffic struct {
	In  string `json:"in"`  // Ex: "2.3 MB/s"
	Out string `json:"out"` // Ex: "0.8 MB/s"
}

type ZPoolVdev struct {
	Name    string   `json:"name"`
	Status  string   `json:"status"` // ONLINE, DEGRADED
	Devices []string `json:"devices"`
}

type ZFSConfig struct {
	PoolName    string      `json:"poolName"`
	PoolStatus  string      `json:"poolStatus"` // ONLINE, DEGRADED
	DataVdevs   []ZPoolVdev `json:"dataVdevs"`
	CacheVdev   *ZPoolVdev  `json:"cacheVdev,omitempty"` // Pointeur pour gérer l'absence de L2ARC
}

type ARCCache struct {
	ARCSize       string  `json:"arcSize"`       // Ex: "64.5 GB"
	ARCTargetSize string  `json:"arcTargetSize"` // Ex: "65.0 GB"
	ARCHitRate    string  `json:"arcHitRate"`    // Ex: "97%"
	ARCHitRateNum float64 `json:"arcHitRateNum"` // Pour le graphique en anneau
	L2ARCSize     string  `json:"l2arcSize"`     // Ex: "480.1 GB"
	L2ARCHitRate  string  `json:"l2arcHitRate"`  // Ex: "68.2 %"
}

type SystemInfo struct {
	OS     string `json:"os"`
	Kernel string `json:"kernel"`
	CPU    string `json:"cpu"`
	Uptime string `json:"uptime"` // Ex: "7d 14h 33m"
}

type DockerInfo struct {
	Containers int `json:"containers"` // Ex: 22
	Images     int `json:"images"`     // Ex: 128
	Volumes    int `json:"volumes"`    // Ex: 35
}

type StreamingInfo struct {
	Films      int `json:"films"`      // Ex: 1923
	Series     int `json:"series"`     // Ex: 812
	Animes     int `json:"animes"`     // Ex: 350
	Playing    int `json:"playing"`    // Ex: 2 (Lectures en cours)
	Transcoding int `json:"transcoding"` // Ex: 1 (Transcodage actif)
}

// GlobalMetrics regroupe toutes les métriques que nous enverrons au frontend
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
		// On retourne une chaîne vide et l'erreur si la commande a échoué.
		return "", fmt.Errorf("erreur lors de l'exécution de '%s %v': %w", name, args, err)
	}
	// On retourne la sortie de la commande (nettoyée des espaces superflus) et pas d'erreur.
	return strings.TrimSpace(string(output)), nil
}

func getStreamingInfo() StreamingInfo {
	// Fonction interne pour compter les éléments dans un dossier
	countItemsInDir := func(path string) int {
		if path == "" {
			return 0 // Si la variable d'env n'est pas définie, on retourne 0
		}
		
		// os.ReadDir lit le contenu d'un dossier
		entries, err := os.ReadDir(path)
		if err != nil {
			log.Printf("Erreur de lecture du dossier %s: %v", path, err)
			return 0 // En cas d'erreur (ex: le dossier n'existe pas), on retourne 0
		}
		return len(entries)
	}

	// On lit les variables d'environnement
	filmsPath := os.Getenv("PATH_FILMS")
	seriesPath := os.Getenv("PATH_SERIES")
	animesPath := os.Getenv("PATH_ANIMES")

	// On utilise notre fonction interne pour compter les éléments
	// Les stats de lecture/transcodage restent mockées pour l'instant
	return StreamingInfo{
		Films:      countItemsInDir(filmsPath),
		Series:     countItemsInDir(seriesPath),
		Animes:     countItemsInDir(animesPath),
		Playing:    2,  // A remplacer plus tard par un appel à l'API de Plex/Jellyfin
		Transcoding: 1,  // Idem
	}
}

func getRAMInfo() RAMInfo {
	// La commande `free -m` donne l'usage de la RAM en mégaoctets.
	out, err := runCommand("free", "-m")
	if err != nil {
		log.Printf("Erreur getRAMInfo: %v", err)
		return RAMInfo{} // Retourne une struct vide en cas d'erreur
	}

	// On parse la sortie de la commande. C'est un peu artisanal mais ça fonctionne.
	lines := strings.Split(out, "\n")
	if len(lines) < 2 {
		return RAMInfo{}
	}
	
	// La ligne qui nous intéresse est la deuxième (index 1)
	fields := strings.Fields(lines[1]) // "Fields" sépare par des espaces
	if len(fields) < 4 {
		return RAMInfo{}
	}

	// On convertit les champs texte en nombres
	total, _ := strconv.ParseFloat(fields[1], 64) // total est le 2ème champ
	used, _ := strconv.ParseFloat(fields[2], 64) // used est le 3ème champ
	
	percent := 0.0
	if total > 0 {
		percent = (used / total) * 100
	}
	
	return RAMInfo{
		Used:    fmt.Sprintf("%.1f GB", used/1024),
		Total:   fmt.Sprintf("%.1f GB", total/1024),
		Percent: fmt.Sprintf("%.1f%%", percent),
		PercentNum: percent,
	}
}


func getSystemInfo() SystemInfo {
	// On exécute plusieurs commandes pour récupérer les infos système.
	osOut, _ := runCommand("cat", "/etc/os-release")
	kernelOut, _ := runCommand("uname", "-r")
	cpuOut, _ := runCommand("lscpu")
	uptimeOut, _ := runCommand("uptime", "-p")

	// Parsing simple pour extraire les infos utiles
	var osName, cpuModel string

	for _, line := range strings.Split(osOut, "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			osName = strings.Trim(strings.Split(line, "=")[1], `"`)
		}
	}

	for _, line := range strings.Split(cpuOut, "\n") {
		if strings.HasPrefix(line, "Model name:") {
			cpuModel = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}
	
	return SystemInfo{
		OS:     osName,
		Kernel: kernelOut,
		CPU:    cpuModel,
		Uptime: strings.TrimPrefix(uptimeOut, "up "),
	}
}

// On crée un squelette pour notre fonction principale de collecte.
// Pour l'instant, elle ne remplit que la RAM et les infos système.
func collectAllMetrics() GlobalMetrics {
	// On crée une structure avec les données mockées qu'on avait avant
	metrics := GlobalMetrics{
		RAM: getRAMInfo(),      // <--- ON UTILISE NOTRE NOUVELLE FONCTION
		System: getSystemInfo(), // <--- ON UTILISE NOTRE NOUVELLE FONCTION
		Streaming: getStreamingInfo(),
		
		// Les autres données sont encore mockées pour l'instant
		CPU: CPUInfo{Usage: "72%", Temp: "51.2°C", TempDeg: 51.2},
		Disk: DiskInfo{Total: "8 TB", Used: "4.5 TB", Free: "3.5 TB", Percent: "56.2%", PercentNum: 56.2},
		Net: NetTraffic{In: "5.7 MB/s", Out: "1.2 MB/s"},
		ZFSConfig: ZFSConfig{PoolName: "rpool", PoolStatus: "ONLINE", DataVdevs: []ZPoolVdev{{Name: "raidz1-0", Status: "ONLINE", Devices: []string{"disk-01", "disk-02", "disk-03", "disk-04"}},{Name: "raidz1-1", Status: "ONLINE", Devices: []string{"disk-05", "disk-06", "disk-07", "disk-08"}}}, CacheVdev: &ZPoolVdev{Name: "L2ARC", Status: "ONLINE", Devices: []string{"nvme-Samsung-970-Evo"}}},
		ARCCache: ARCCache{ARCSize: "66.1 GB", ARCTargetSize: "65.0 GB", ARCHitRate: "98.1%", ARCHitRateNum: 98.1, L2ARCSize: "485.0 GB", L2ARCHitRate: "70.5 %"},
		Docker: DockerInfo{Containers: 25, Images: 132, Volumes: 38},
		Streaming: StreamingInfo{Films: 1950, Series: 820, Animes: 360, Playing: 3, Transcoding: 1},
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

	fmt.Println("Nouveau client connecté au WebSocket")

	for {
		// ON REMPLACE TOUT LE BLOC DE MOCK PAR CET APPEL
		metrics := collectAllMetrics()

		// On encode nos structs en JSON
		jsonMessage, err := json.Marshal(metrics)
		if err != nil {
			log.Printf("Erreur d'encodage JSON: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = ws.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			log.Printf("Erreur d'envoi: %v", err)
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

	fmt.Printf("Serveur démarré. Rendez-vous sur http://localhost:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}