package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port      int
	Env       string
	DBURL     string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		Port:      getEnvAsInt("PORT", 8080),
		Env:       getEnv("ENV", "development"),
		DBURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ukoni?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "super-secret-key"),
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
