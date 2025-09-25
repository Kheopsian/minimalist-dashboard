package main

import (
	"fmt"
	"log"
	"net/http"

	"minimalist-dashboard/internal/config"
	"minimalist-dashboard/internal/handlers"
	"minimalist-dashboard/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

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
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}