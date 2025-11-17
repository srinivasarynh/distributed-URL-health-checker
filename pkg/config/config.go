package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	URLs          []string
	Port          string
	CheckInterval time.Duration
	Timeout       time.Duration
}

func Load() *Config {
	cfg := &Config{
		Port:          getEnv("PORT", "8080"),
		CheckInterval: getDuration("CHECK_INTERVAL", 10*time.Second),
		Timeout:       getDuration("TIMEOUT", 5*time.Second),
	}

	urlsEnv := getEnv("URLS", "https://google.com,http://github.com,http://golang.org")
	cfg.URLs = strings.Split(urlsEnv, ",")

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}

	return fallback
}
