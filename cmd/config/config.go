package config

import (
	"os"
	"strconv"
	"time"
)

//создание конфига для сервера
type ServerConfig struct {
	Addr          string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
}

//создание конфига для агента
type AgentConfig struct {
	Addr           string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

// инициализация конфига для сервера
func NewConfigServer(flagAddress, flagStoreInterval, flagStoreFile, flagRestore string) *ServerConfig {
	return &ServerConfig{
		//проверка флагов и облачных переменных, приоритет облачным переменным, если не дефолтные значения
		Addr:          getEnv("ADDRESS", flagAddress, "localhost:8080"),
		StoreInterval: getEnvTime("STORE_INTERVAL", flagStoreInterval, "300s"),
		StoreFile:     getEnv("STORE_FILE", flagStoreFile, "/tmp/devops-metrics-db.json"),
		Restore:       getEnvAsBool("RESTORE", flagRestore, "true"),
	}
}
func NewConfigAgent(flagAddress, flagPollInterval, flagReportInterval string) *AgentConfig {
	return &AgentConfig{
		//проверка флагов и облачных переменных, приоритет облачным переменным, если не дефолтные значения
		Addr:           getEnv("ADDRESS", flagAddress, "localhost:8080"),
		PollInterval:   getEnvTime("POLL_INTERVAL", flagPollInterval, "2s"),
		ReportInterval: getEnvTime("REPORT_INTERVAL", flagReportInterval, "10s"),
	}
}

//сбор значений из облачных переменных и флагов в формате string
func getEnv(key, flag, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if flag != defaultVal {
		return flag
	}
	return defaultVal
}

//сбор значений из облачных переменных и флагов в формате time.Duration
func getEnvTime(key, flag, defaultVal string) time.Duration {
	valueStr := getEnv(key, flag, defaultVal)
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		value, _ = time.ParseDuration(defaultVal)
		return value
	}
	return value
}

//сбор значений из облачных переменных и флагов в формате bool
func getEnvAsBool(key, flag string, defaultVal string) bool {
	valStr := getEnv(key, flag, defaultVal)
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return true
}
