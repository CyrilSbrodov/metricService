package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

// AgentConfig структура конфига для агента.
type AgentConfig struct {
	Addr           string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	Hash           string        `env:"KEY"`
}

var cfgAgent AgentConfig

// AgentConfigInit инициализая конфига.
func AgentConfigInit() AgentConfig {
	flag.StringVar(&cfgAgent.Addr, "a", "http://localhost:8080", "ADDRESS")
	flag.DurationVar(&cfgAgent.PollInterval, "p", time.Duration(2)*time.Second, "update interval")
	flag.DurationVar(&cfgAgent.ReportInterval, "r", time.Duration(10)*time.Second, "upload interval")
	flag.StringVar(&cfgAgent.Hash, "k", "", "hash")

	flag.Parse()
	if err := env.Parse(&cfgAgent); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cfgAgent
}
