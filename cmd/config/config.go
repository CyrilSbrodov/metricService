package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

//объявление флагов
func ServerFlagsInit() (flagAddress, flagRestore, flagStoreInterval, flagStoreFile, flagHash, flagDatabase *string) {
	//присвоение значений флагам
	flagAddress = flag.String("a", "localhost:8080", "address of service")
	flagRestore = flag.String("r", "true", "restore from file")
	flagStoreInterval = flag.String("i", "300s", "upload interval")
	flagStoreFile = flag.String("f", "/tmp/devops-metrics-db.json", "name of file")
	flagHash = flag.String("k", "КЛЮЧ", "hash")
	flagDatabase = flag.String("d", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "name of database")
	return
	//
}

func AgentFlagsInit() (flagAddress, flagPollInterval, flagReportInterval, flagHash *string) {
	//присвоение значений флагам
	flagAddress = flag.String("a", "localhost:8080", "address of service")
	flagPollInterval = flag.String("p", "2s", "update interval")
	flagReportInterval = flag.String("r", "10s", "upload interval to server")
	flagHash = flag.String("k", "КЛЮЧ", "hash")
	return

}

//создание конфига для сервера
type ServerConfig struct {
	Addr          string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	Hash          string
	DatabaseDSN   string
}

//создание конфига для агента
type AgentConfig struct {
	Addr           string
	PollInterval   time.Duration
	ReportInterval time.Duration
	Hash           string
}

// инициализация конфига для сервера
func NewConfigServer(flagAddress, flagStoreInterval, flagStoreFile, flagRestore, flagHash, flagDatabase string) *ServerConfig {
	return &ServerConfig{
		//проверка флагов и облачных переменных, приоритет облачным переменным, если не дефолтные значения
		Addr:          getEnv("ADDRESS", flagAddress, "localhost:8080"),
		StoreInterval: getEnvTime("STORE_INTERVAL", flagStoreInterval, "300s"),
		StoreFile:     getEnv("STORE_FILE", flagStoreFile, "/tmp/devops-metrics-db.json"),
		Restore:       getEnvAsBool("RESTORE", flagRestore, "true"),
		Hash:          getEnv("KEY", flagHash, "КЛЮЧ"),
		DatabaseDSN:   getEnv("DATABASE_DSN", flagDatabase, ""),
	}
}
func NewConfigAgent(flagAddress, flagPollInterval, flagReportInterval, flagHash string) *AgentConfig {
	return &AgentConfig{
		//проверка флагов и облачных переменных, приоритет облачным переменным, если не дефолтные значения
		Addr:           getEnv("ADDRESS", flagAddress, "localhost:8080"),
		PollInterval:   getEnvTime("POLL_INTERVAL", flagPollInterval, "2s"),
		ReportInterval: getEnvTime("REPORT_INTERVAL", flagReportInterval, "10s"),
		Hash:           getEnv("KEY", flagHash, "КЛЮЧ"),
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
