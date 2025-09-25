package config

import "os"

// Config contient toute la configuration de l'application
type Config struct {
	WebUIPort     string
	PathFilms     string
	PathSeries    string
	PathAnimes    string
	NetInterface  string
	PlexURL       string
	PlexToken     string
}

// Load charge la configuration depuis les variables d'environnement
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

// getEnvOrDefault retourne la valeur de la variable d'environnement ou la valeur par défaut
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetMonitorPath retourne le chemin à surveiller pour le disque (PATH_FILMS ou racine)
func (c *Config) GetMonitorPath() string {
	if c.PathFilms != "" {
		return c.PathFilms
	}
	return "/"
}