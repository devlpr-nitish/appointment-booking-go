package config

import (
	"os"
)

type Config struct {
	DBDriver string
	DBUrl    string
	AppPort  string
}

func LoadConfig() *Config {
	return &Config{
		DBDriver: getEnv("DB_DRIVER", "postgres"), // Default to postgres
		DBUrl:    getEnv("DB_URL", ""),
		AppPort:  getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
