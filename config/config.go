// Package config loads all application configuration from environment variables.
// The .env file is loaded by the caller (main.go) before this package is used.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config is the top-level application configuration.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// AppConfig holds HTTP server settings.
type AppConfig struct {
	Port    string
	GinMode string // "debug" | "release"
}

// DatabaseConfig holds MySQL connection settings.
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DSN returns the MySQL Data Source Name string.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name,
	)
}

// JWTConfig holds JWT signing settings.
type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

// Load reads configuration from environment variables and returns a populated Config.
// It returns an error if any required variable is missing.
func Load() (*Config, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("config: JWT_SECRET environment variable is required")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	// DB_PASSWORD is intentionally allowed to be empty in local-only setups,
	// but we warn loudly in production (GIN_MODE=release).
	if dbPassword == "" && os.Getenv("GIN_MODE") == "release" {
		return nil, fmt.Errorf("config: DB_PASSWORD must not be empty in release mode")
	}

	maxOpen := getEnvInt("DB_MAX_OPEN_CONNS", 25)
	maxIdle := getEnvInt("DB_MAX_IDLE_CONNS", 10)
	lifetimeMin := getEnvInt("DB_CONN_MAX_LIFETIME_MIN", 5)
	jwtExpHours := getEnvInt("JWT_EXPIRATION_HOURS", 72)

	cfg := &Config{
		App: AppConfig{
			Port:    getEnvStr("APP_PORT", "8080"),
			GinMode: getEnvStr("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:            getEnvStr("DB_HOST", "localhost"),
			Port:            getEnvStr("DB_PORT", "3306"),
			User:            getEnvStr("DB_USER", "root"),
			Password:        dbPassword,
			Name:            getEnvStr("DB_NAME", "clean_anti_db"),
			MaxOpenConns:    maxOpen,
			MaxIdleConns:    maxIdle,
			ConnMaxLifetime: time.Duration(lifetimeMin) * time.Minute,
		},
		JWT: JWTConfig{
			Secret:          jwtSecret,
			ExpirationHours: jwtExpHours,
		},
	}

	return cfg, nil
}

// --- helpers ---

func getEnvStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
