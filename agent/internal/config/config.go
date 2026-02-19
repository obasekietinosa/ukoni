package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         int
	APIBaseURL   string
	CORSAllowed  string
	Environment  string
	DBURL        string
}

func Load() *Config {
	return &Config{
		Port:         getEnvAsInt("PORT", 8081),
		APIBaseURL:   getEnv("API_BASE_URL", "http://localhost:8080"),
		CORSAllowed:  getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:5173"),
		Environment:  getEnv("ENV", "development"),
		DBURL:        getEnv("DATABASE_URL", "postgres://etin:etin@localhost:5432/ukoni?sslmode=disable"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
