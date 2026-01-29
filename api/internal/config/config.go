package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port               int
	Env                string
	DBURL              string
	JWTSecret          string
	CorsAllowedOrigins []string
}

func Load() *Config {
	return &Config{
		Port:               getEnvAsInt("PORT", 8080),
		Env:                getEnv("ENV", "development"),
		DBURL:              getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ukoni?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "super-secret-key"),
		CorsAllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}, ","),
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

func getEnvAsSlice(key string, defaultVal []string, sep string) []string {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	parts := strings.Split(valStr, sep)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
