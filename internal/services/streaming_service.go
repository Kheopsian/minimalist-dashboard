package services

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"minimalist-dashboard/internal/config"
	"minimalist-dashboard/internal/models"
)

// StreamingService gère les informations de streaming (Plex)
type StreamingService struct {
	config *config.Config
}

// NewStreamingService crée une nouvelle instance du service streaming
func NewStreamingService(cfg *config.Config) *StreamingService {
	return &StreamingService{config: cfg}
}

// GetStreamingInfo récupère les informations sur les services de streaming
func (s *StreamingService) GetStreamingInfo() models.StreamingInfo {
	// -- Partie 1 : Compter les fichiers --
	countItemsInDir := func(path string) int {
		if path == "" {
			return 0
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return 0
		}
		return len(entries)
	}
	films := countItemsInDir(s.config.PathFilms)
	series := countItemsInDir(s.config.PathSeries)
	animes := countItemsInDir(s.config.PathAnimes)

	// -- Partie 2 : Interroger l'API Plex --
	playing := 0
	transcoding := 0

	// Si l'URL ou le Token ne sont pas définis, on saute cette partie
	if s.config.PlexURL != "" && s.config.PlexToken != "" {
		// On crée un client HTTP avec un timeout de 2 secondes
		client := &http.Client{Timeout: 2 * time.Second}

		// On prépare la requête GET vers l'endpoint /status/sessions
		req, err := http.NewRequest("GET", s.config.PlexURL+"/status/sessions", nil)
		if err == nil {
			// On ajoute les en-têtes nécessaires pour l'authentification et le format
			req.Header.Add("Accept", "application/json")
			req.Header.Add("X-Plex-Token", s.config.PlexToken)

			// On exécute la requête
			resp, err := client.Do(req)
			if err == nil {
				defer resp.Body.Close()

				// On décode la réponse JSON dans nos structs
				var sessionResponse models.PlexSessionResponse
				err = json.NewDecoder(resp.Body).Decode(&sessionResponse)
				if err == nil {
					// On a réussi ! On met à jour nos compteurs.
					playing = sessionResponse.MediaContainer.Size
					for _, media := range sessionResponse.MediaContainer.Metadata {
						if len(media.TranscodeSession) > 0 {
							transcoding++
						}
					}
				}
			}
		}
		if err != nil {
			log.Printf("Erreur lors de l'appel à l'API Plex: %v", err)
		}
	}

	return models.StreamingInfo{
		Films:       films,
		Series:      series,
		Animes:      animes,
		Playing:     playing,
		Transcoding: transcoding,
	}
}