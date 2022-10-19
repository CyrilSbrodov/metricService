package config

import (
	"os"
	"time"
)

type Config struct {
	Addr           string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewConfig() *Config {
	return &Config{
		Addr:           getEnv("ADDRESS", "localhost:8080"),
		PollInterval:   getEnvTime("POLL_INTERVAL", 2*time.Second),
		ReportInterval: getEnvTime("REPORT_INTERVAL", 10*time.Second),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvTime(key string, defaultVal time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultVal
	}
	return value
}
