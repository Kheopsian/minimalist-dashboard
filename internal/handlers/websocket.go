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

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	config         *config.Config
	metricsService *services.MetricsService
}

// NewWebSocketHandler creates a new WebSocket handler instance
func NewWebSocketHandler(cfg *config.Config, metricsService *services.MetricsService) *WebSocketHandler {
	return &WebSocketHandler{
		config:         cfg,
		metricsService: metricsService,
	}
}

// HandleConnections handles incoming WebSocket connections
func (h *WebSocketHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	fmt.Println("New client connected to WebSocket")

	// Initialize variables for delta calculations
	cpuService := services.NewCPUService()
	storageService := services.NewStorageService(h.config)
	
	var prevCPUTimes models.CPUTimes
	var prevNetCounters models.NetCounters
	prevCPUTimes, _ = cpuService.GetCPUTimes()
	prevNetCounters, _ = storageService.GetNetCounters(h.config.NetInterface)
	prevTime := time.Now()

	for {
		time.Sleep(2 * time.Second)

		// Collect real-time metrics
		metrics, newCPUTimes, newNetCounters, newTime := h.metricsService.CollectRealTimeMetrics(
			prevCPUTimes,
			prevNetCounters,
			prevTime,
		)

		// Update variables for next iteration
		prevCPUTimes = newCPUTimes
		prevNetCounters = newNetCounters
		prevTime = newTime

		// Send data to client
		jsonMessage, _ := json.Marshal(metrics)
		err = ws.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			log.Printf("Send error: %v", err)
			break
		}
	}
}