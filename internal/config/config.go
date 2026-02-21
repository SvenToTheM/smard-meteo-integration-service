package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Port         string
	SmardBaseURL string
	MeteoBaseURL string
	DefaultLat   float64
	DefaultLon   float64
}

// LoadConfig reads configuration from environment variables or falls back to defaults
func LoadConfig() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		SmardBaseURL: getEnv("SMARD_BASE_URL", "https://www.smard.de/app"),
		MeteoBaseURL: getEnv("METEO_BASE_URL", "https://archive-api.open-meteo.com/v1/archive"), // For past dates
		DefaultLat:   getEnvAsFloat("DEFAULT_LAT", 51.1),
		DefaultLon:   getEnvAsFloat("DEFAULT_LON", 10.4),
	}
}

// Helper to read an environment variable into a string or return a fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper to read an environment variable into a float64 or return a fallback
func getEnvAsFloat(key string, fallback float64) float64 {
	strValue := getEnv(key, "")
	if value, err := strconv.ParseFloat(strValue, 64); err == nil {
		return value
	}
	return fallback
}
