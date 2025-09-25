package config

import "os"

// Config contains all the application configuration
type Config struct {
	WebUIPort     string
	PathFilms     string
	PathSeries    string
	PathAnimes    string
	NetInterface  string
	PlexURL       string
	PlexToken     string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		WebUIPort:    getEnvOrDefault("WEBUI_PORT", "8080"),
		PathFilms:    getEnvOrDefault("PATH_FILMS", ""),
		PathSeries:   getEnvOrDefault("PATH_SERIES", ""),
		PathAnimes:   getEnvOrDefault("PATH_ANIMES", ""),
		NetInterface: getEnvOrDefault("NET_INTERFACE", "eth0"),
		PlexURL:      getEnvOrDefault("PLEX_URL", ""),
		PlexToken:    getEnvOrDefault("PLEX_TOKEN", ""),
	}
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