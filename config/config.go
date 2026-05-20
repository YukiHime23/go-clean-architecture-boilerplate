package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	App AppConfig
	DB  DBConfig
	JWT JWTConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	DSN      string
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

func Load() *Config {
	expireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "72"))
	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		DB: newDBConfig(),
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "change-me-in-production"),
			ExpireHours: expireHours,
		},
	}
}

func newDBConfig() DBConfig {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "password")
	name := getEnv("DB_NAME", "go_clean_db")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, name)
	return DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     name,
		DSN:      dsn,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
