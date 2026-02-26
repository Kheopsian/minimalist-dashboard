package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config contains all the application configuration
type Config struct {
	WebUIPort    string `json:"WebUIPort"`
	PathFilms    string `json:"PathFilms"`
	PathSeries   string `json:"PathSeries"`
	PathAnimes   string `json:"PathAnimes"`
	NetInterface string `json:"NetInterface"`
	PlexURL      string `json:"PlexURL"`
	PlexToken    string `json:"PlexToken"`
}

// Load loads configuration from environment variables (fallback/legacy)
func Load() *Config {
	return &Config{
		WebUIPort:    getEnvOrDefault("WEBUI_PORT", "9254"),
		PathFilms:    getEnvOrDefault("PATH_FILMS", ""),
		PathSeries:   getEnvOrDefault("PATH_SERIES", ""),
		PathAnimes:   getEnvOrDefault("PATH_ANIMES", ""),
		NetInterface: getEnvOrDefault("NET_INTERFACE", "eth0"),
		PlexURL:      getEnvOrDefault("PLEX_URL", ""),
		PlexToken:    getEnvOrDefault("PLEX_TOKEN", ""),
	}
}

// LoadFromFile loads configuration from a JSON file
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	return &cfg, nil
}

// getEnvOrDefault returns the environment variable value or the default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetMonitorPath returns the path to monitor for disk usage (PATH_FILMS or root)
func (c *Config) GetMonitorPath() string {
	if c.PathFilms != "" {
		return c.PathFilms
	}
	return "/"
}