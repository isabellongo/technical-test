package config

import (
	"os"
)

// Config holds application configuration from environment.
type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// API
	APIPort string
	APIKey  string
}

// Load reads configuration from environment variables.
func Load() *Config {
	return &Config{
		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnv("POSTGRES_PORT", "5432"),
		DBUser:     getEnv("POSTGRES_USER", "postgres"),
		DBPassword: getEnv("POSTGRES_PASSWORD", "changeme"),
		DBName:     getEnv("POSTGRES_DB", "driva"),
		APIPort:    getEnv("API_PORT", "3000"),
		// API_KEY do .env; fallback para chave de teste do desafio
		APIKey: getEnv("API_KEY", "driva_test_key_abc123xyz789"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
