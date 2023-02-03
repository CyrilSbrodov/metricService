/*
	Package config создание конфиг файла для сервера
*/
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

// ServerConfig структура конфига для сервера.
type ServerConfig struct {
	Addr          string        `json:"address" env:"ADDRESS"`
	Config        string        `env:"CONFIG"`
	StoreFile     string        `json:"store_file" env:"STORE_FILE"`
	Hash          string        `env:"KEY"`
	DatabaseDSN   string        `json:"database_dsn" env:"DATABASE_DSN"`
	CryptoPROKey  string        `json:"crypto_key" env:"CRYPTO_KEY"`
	Restore       bool          `json:"restore" env:"RESTORE"`
	StoreInterval time.Duration `json:"store_interval" env:"STORE_INTERVAL"`
}

var cfgSrv ServerConfig

// ServerConfigInit инициализация конфига.
func ServerConfigInit() ServerConfig {
	cfg := parseFromServerConfigFile()
	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "ADDRESS")
	flag.DurationVar(&cfg.StoreInterval, "i", time.Duration(300)*time.Second, "STORE_INTERVAL")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "STORE_FILE")
	flag.BoolVar(&cfg.Restore, "r", true, "RESTORE")
	flag.StringVar(&cfg.Hash, "k", "", "KEY")
	flag.StringVar(&cfg.DatabaseDSN, "d", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "DATABASE_DSN")
	flag.StringVar(&cfg.CryptoPROKey, "crypto-key", "", "path to file")
	// ../../internal/crypto/privateKeyPEM
	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cfg
}

//parseFromServerConfigFile загрузка конфига из файла
func parseFromServerConfigFile() ServerConfig {
	flag.StringVar(&cfgSrv.Config, "c/-config", "", "path to config file")
	flag.Parse()
	if err := env.Parse(&cfgSrv); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if cfgSrv.Config != "" {
		configFile, err := os.Open(cfgSrv.Config)
		defer configFile.Close()
		if err != nil {
			fmt.Println(err.Error())
			//TODO
		}
		jsonParser := json.NewDecoder(configFile)
		jsonParser.Decode(&cfgSrv)
		return cfgSrv
	}
	return cfgSrv
}
