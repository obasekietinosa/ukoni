package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port               int
	APIBaseURL         string
	CorsAllowedOrigins []string
	Environment        string
	DBURL              string
}

func Load() *Config {
	return &Config{
		Port:               getEnvAsInt("PORT", 8081),
		APIBaseURL:         getEnv("API_BASE_URL", "http://localhost:8080"),
		CorsAllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:5173"}, ","),
		Environment:        getEnv("ENV", "development"),
		DBURL:              getEnv("DATABASE_URL", "postgres://etin:etin@localhost:5432/ukoni?sslmode=disable"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsSlice(key string, defaultVal []string, sep string) []string {
	valStr := getEnv(key, "")
	if valStr == "" {
		// Fallback to singular just in case it was used
		if key == "CORS_ALLOWED_ORIGINS" {
			if single := getEnv("CORS_ALLOWED_ORIGIN", ""); single != "" {
				valStr = single
			} else {
				return defaultVal
			}
		} else {
			return defaultVal
		}
	}
	parts := strings.Split(valStr, sep)
	for i := range parts {
		p := strings.TrimSpace(parts[i])
		p = strings.Trim(p, "\"'")
		p = strings.TrimRight(p, "/")
		parts[i] = p
	}
	return parts
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
