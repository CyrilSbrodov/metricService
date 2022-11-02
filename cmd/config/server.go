package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Addr          string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Hash          string        `env:"KEY"`
	DatabaseDSN   string        `env:"DATABASE_DSN"`
}

var cfgSrv ServerConfig

func ServerConfigInit() ServerConfig {
	flag.StringVar(&cfgSrv.Addr, "a", "localhost:8080", "server address")
	flag.DurationVar(&cfgSrv.StoreInterval, "i", time.Duration(300)*time.Second, "upload interval")
	flag.StringVar(&cfgSrv.StoreFile, "f", "/tmp/devops-metrics-db.json", "file name")
	flag.BoolVar(&cfgSrv.Restore, "r", true, "restore from file")
	flag.StringVar(&cfgSrv.Hash, "k", "", "hash")
	flag.StringVar(&cfgSrv.DatabaseDSN, "d", "", "database addr")

	flag.Parse()
	if err := env.Parse(&cfgSrv); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cfgSrv
}
