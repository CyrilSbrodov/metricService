package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr           string
	PollInterval   time.Duration
	ReportInterval time.Duration
	StoreInterval  time.Duration
	StoreFile      string
	Restore        bool
}

func NewConfig() *Config {
	return &Config{
		Addr:           getEnv("ADDRESS", "localhost:8080"),
		PollInterval:   getEnvTime("POLL_INTERVAL", 2*time.Second),
		ReportInterval: getEnvTime("REPORT_INTERVAL", 10*time.Second),
		StoreInterval:  getEnvTime("STORE_INTERVAL", 300*time.Second),
		StoreFile:      getEnv("STORE_FILE", "/tmp/devops-metrics-db.json"),
		Restore:        getEnvAsBool("RESTORE", true),
	}
}

//"/tmp/devops-metrics-db.json"
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

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
