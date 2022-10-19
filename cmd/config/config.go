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
		Addr:          getEnv(flagAddress, "localhost:8080"),
		StoreInterval: getEnvTime(flagStoreInterval, 300*time.Second),
		StoreFile:     getEnv(flagStoreFile, "/tmp/devops-metrics-db.json"),
		Restore:       getEnvAsBool(flagRestore, true),
	}
}
func NewConfigAgent(flagAddress, flagPollInterval, flagReportInterval string) *AgentConfig {
	return &AgentConfig{
		Addr:           getEnv(flagAddress, "localhost:8080"),
		PollInterval:   getEnvTime(flagPollInterval, 2*time.Second),
		ReportInterval: getEnvTime(flagReportInterval, 10*time.Second),
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
