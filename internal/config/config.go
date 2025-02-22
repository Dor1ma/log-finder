package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	LogDir          string
	ServerPort      string
	CacheTTL        time.Duration
	MaxOpenFiles    int
	FileCacheTTL    time.Duration
	RateLimit       int
	RefreshInterval time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system default variables")
	}

	return &Config{
		LogDir:          getEnv("LOG_DIR", "/tmp/log"),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		CacheTTL:        getEnvAsDuration("CACHE_TTL", 5*time.Minute),
		MaxOpenFiles:    getEnvAsInt("MAX_OPEN_FILES", 20),
		FileCacheTTL:    getEnvAsDuration("FILE_CACHE_TTL", 10*time.Minute),
		RateLimit:       getEnvAsInt("RATE_LIMIT", 100),
		RefreshInterval: getEnvAsDuration("REFRESH_INERVAL", 60*time.Minute),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		duration, err := time.ParseDuration(value)
		if err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
