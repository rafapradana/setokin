// Package config provides configuration management for the Setokin API.
// It loads configuration from environment variables.
package config

import (
	"os"
	"time"
)

// Config holds the complete application configuration.
type Config struct {
	App   AppConfig
	DB    DBConfig
	JWT   JWTConfig
	MinIO MinIOConfig
	CORS  CORSConfig
}

// AppConfig holds application-level settings.
type AppConfig struct {
	Name string
	Env  string
	Port string
}

// DBConfig holds database connection settings.
type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

// JWTConfig holds JWT authentication settings.
type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

// MinIOConfig holds MinIO object storage settings.
type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// CORSConfig holds CORS settings.
type CORSConfig struct {
	AllowedOrigins string
}

// Load reads configuration from environment variables and returns a Config.
func Load() *Config {
	return &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "setokin"),
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("API_PORT", "8080"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "setokin"),
			User:     getEnv("DB_USER", "setokin"),
			Password: getEnv("DB_PASSWORD", "setokin_password"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "change-me-in-production"),
			AccessExpiry:  parseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m")),
			RefreshExpiry: parseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h")),
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "setokin"),
			UseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
		},
	}
}

// getEnv reads an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// parseDuration parses a duration string, falling back to 15 minutes on error.
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}
