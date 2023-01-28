/*
	Package config создание конфиг файла для сервера
*/
package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

// ServerConfig структура конфига для сервера.
type ServerConfig struct {
	Addr          string        `env:"ADDRESS"`
	StoreFile     string        `env:"STORE_FILE"`
	Hash          string        `env:"KEY"`
	DatabaseDSN   string        `env:"DATABASE_DSN"`
	Restore       bool          `env:"RESTORE"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
}

var cfgSrv ServerConfig

// ServerConfigInit инициализация конфига.
func ServerConfigInit() ServerConfig {
	flag.StringVar(&cfgSrv.Addr, "a", "localhost:8080", "ADDRESS")
	flag.DurationVar(&cfgSrv.StoreInterval, "i", time.Duration(300)*time.Second, "STORE_INTERVAL")
	flag.StringVar(&cfgSrv.StoreFile, "f", "/tmp/devops-metrics-db.json", "STORE_FILE")
	flag.BoolVar(&cfgSrv.Restore, "r", true, "RESTORE")
	flag.StringVar(&cfgSrv.Hash, "k", "", "KEY")
	flag.StringVar(&cfgSrv.DatabaseDSN, "d", "", "DATABASE_DSN")
	flag.Parse()
	if err := env.Parse(&cfgSrv); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cfgSrv
}
