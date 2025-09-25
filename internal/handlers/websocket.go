package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"minimalist-dashboard/internal/config"
	"minimalist-dashboard/internal/models"
	"minimalist-dashboard/internal/services"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocketHandler gère les connexions WebSocket
type WebSocketHandler struct {
	config         *config.Config
	metricsService *services.MetricsService
}

// NewWebSocketHandler crée une nouvelle instance du handler WebSocket
func NewWebSocketHandler(cfg *config.Config, metricsService *services.MetricsService) *WebSocketHandler {
	return &WebSocketHandler{
		config:         cfg,
		metricsService: metricsService,
	}
}

// HandleConnections gère les connexions WebSocket entrantes
func (h *WebSocketHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	fmt.Println("Nouveau client connecté au WebSocket")

	// Initialisation des variables pour le calcul des deltas
	cpuService := services.NewCPUService()
	storageService := services.NewStorageService(h.config)
	
	var prevCPUTimes models.CPUTimes
	var prevNetCounters models.NetCounters
	prevCPUTimes, _ = cpuService.GetCPUTimes()
	prevNetCounters, _ = storageService.GetNetCounters(h.config.NetInterface)
	prevTime := time.Now()

	for {
		time.Sleep(2 * time.Second)

		// Collecte des métriques en temps réel
		metrics, newCPUTimes, newNetCounters, newTime := h.metricsService.CollectRealTimeMetrics(
			prevCPUTimes,
			prevNetCounters,
			prevTime,
		)

		// Mise à jour des variables pour la prochaine itération
		prevCPUTimes = newCPUTimes
		prevNetCounters = newNetCounters
		prevTime = newTime

		// Envoi des données au client
		jsonMessage, _ := json.Marshal(metrics)
		err = ws.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			log.Printf("Erreur d'envoi: %v", err)
			break
		}
	}
}