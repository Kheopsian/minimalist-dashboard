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

// StreamingService manages streaming information (Plex)
type StreamingService struct {
	config *config.Config
}

// NewStreamingService creates a new streaming service instance
func NewStreamingService(cfg *config.Config) *StreamingService {
	return &StreamingService{config: cfg}
}

// GetStreamingInfo retrieves streaming services information
func (s *StreamingService) GetStreamingInfo() models.StreamingInfo {
	// -- Part 1: Count files --
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

	// -- Part 2: Query Plex API --
	playing := 0
	transcoding := 0

	// If URL or Token are not defined, skip this part
	if s.config.PlexURL != "" && s.config.PlexToken != "" {
		// Create an HTTP client with a 2-second timeout
		client := &http.Client{Timeout: 2 * time.Second}

		// Prepare GET request to the /status/sessions endpoint
		req, err := http.NewRequest("GET", s.config.PlexURL+"/status/sessions", nil)
		if err == nil {
			// Add necessary headers for authentication and format
			req.Header.Add("Accept", "application/json")
			req.Header.Add("X-Plex-Token", s.config.PlexToken)

			// Execute the request
			resp, err := client.Do(req)
			if err == nil {
				defer resp.Body.Close()

				// Decode JSON response into our structs
				var sessionResponse models.PlexSessionResponse
				err = json.NewDecoder(resp.Body).Decode(&sessionResponse)
				if err == nil {
					// Success! Update our counters.
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
			log.Printf("Error calling Plex API: %v", err)
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