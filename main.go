package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"minimalist-dashboard/internal/config"
	"minimalist-dashboard/internal/handlers"
	"minimalist-dashboard/internal/services"
)

func main() {
	// Load configuration
	configPath := "data/config.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "config.json"
	}

	cfg, err := config.LoadFromFile(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v\nPlease copy config.example.json to config.json and adjust it.", err)
	}

	// Initialize services
	metricsService := services.NewMetricsService(cfg)

	// Initialize handlers
	wsHandler := handlers.NewWebSocketHandler(cfg, metricsService)

	// Configure routes
	fileServer := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fileServer)
	http.HandleFunc("/ws", wsHandler.HandleConnections)

	// Start server
	listenAddr := ":" + cfg.WebUIPort
	fmt.Printf("Server started. Go to http://localhost:%s\n", cfg.WebUIPort)
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}