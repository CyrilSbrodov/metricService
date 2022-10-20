package config

import (
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	Addr          string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
}

type AgentConfig struct {
	Addr           string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewConfigServer(flagAddress, flagStoreInterval, flagStoreFile, flagRestore string) *ServerConfig {
	return &ServerConfig{
		Addr:          getEnv("ADDRESS", flagAddress, "localhost:8080"),
		StoreInterval: getEnvTime("STORE_INTERVAL", flagStoreInterval, "300s"),
		StoreFile:     getEnv("STORE_FILE", flagStoreFile, "/tmp/devops-metrics-db.json"),
		Restore:       getEnvAsBool("RESTORE", flagRestore, "true"),
	}
}
func NewConfigAgent(flagAddress, flagPollInterval, flagReportInterval string) *AgentConfig {
	return &AgentConfig{
		Addr:           getEnv("ADDRESS", flagAddress, "localhost:8080"),
		PollInterval:   getEnvTime("POLL_INTERVAL", flagPollInterval, "2s"),
		ReportInterval: getEnvTime("REPORT_INTERVAL", flagReportInterval, "10s"),
	}
}

func getEnv(key, flag, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if flag != defaultVal {
		return flag
	}
	return defaultVal
}

func getEnvTime(key, flag, defaultVal string) time.Duration {
	valueStr := getEnv(key, flag, defaultVal)
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return 1
	}
	return value
}

func getEnvAsBool(key, flag string, defaultVal string) bool {
	valStr := getEnv(key, flag, defaultVal)
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return true
}
