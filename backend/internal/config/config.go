package config

import (
	"log"
	"os"

	"dalivim/internal/database"
)

type Config struct {
	Database database.Config
	Server   ServerConfig
}

type ServerConfig struct {
	Port string
	Host string
}

func Load() *Config {
	return &Config{
		Database: database.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5433"),
			User:     getEnv("DB_USER", "p"),
			Password: getEnv("DB_PASSWORD", "p"),
			DBName:   getEnv("DB_NAME", "dalivim"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue != "" {
			log.Printf("Using default for %s: %s", key, defaultValue)
		}
		return defaultValue
	}
	return value
}
