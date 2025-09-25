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
	// Chargement de la configuration
	cfg := config.Load()

	// Initialisation des services
	metricsService := services.NewMetricsService(cfg)

	// Initialisation des handlers
	wsHandler := handlers.NewWebSocketHandler(cfg, metricsService)

	// Configuration des routes
	fileServer := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fileServer)
	http.HandleFunc("/ws", wsHandler.HandleConnections)

	// Démarrage du serveur
	listenAddr := ":" + cfg.WebUIPort
	fmt.Printf("Serveur démarré. Rendez-vous sur http://localhost:%s\n", cfg.WebUIPort)
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}