// Package config provides application configuration loading
// from environment variables (using config.env support via godotenv).
//
// It defines the structure for application and database settings,
// and exposes a LoadConfig function to initialize them.
package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Host holds the server's host and port configuration.
type Host struct {
	ServerHost string
	ServerPort string
}

// Db holds the database connection configuration.
type Db struct {
	Host     string
	Port     string
	User     string
	Password string
	Db       string
	Driver   string
}

// Config combines all app configuration sections.
type Config struct {
	Host Host
	Db   Db
}

// LoadConfig loads configuration from environment variables (with config.env fallback).
//
// It returns a populated Config struct with default values if variables are missing.
func LoadConfig() *Config {
	_ = godotenv.Load("config.env")

	return &Config{
		Host: Host{
			ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),
			ServerPort: getEnv("SERVER_PORT", "8080"),
		},
		Db: Db{
			Host:     getEnv("DB_HOST", "0.0.0.0"),
			User:     getEnv("DB_USER", "admin"),
			Password: getEnv("DB_PASSWORD", "password"),
			Db:       getEnv("DB_NAME", "wrallet"),
			Port:     getEnv("DB_PORT", "5432"),
			Driver:   getEnv("DRIVER", "postgres"),
		},
	}
}

// getEnv returns the value of the environment variable or a default if not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
