package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type AgentConfig struct {
	Addr           string
	PollInterval   time.Duration
	ReportInterval time.Duration
	Hash           string
}

var cfgAgent AgentConfig

func AgentConfigInit() AgentConfig {
	flag.StringVar(&cfgAgent.Addr, "a", "localhost:8080", "server address")
	flag.DurationVar(&cfgAgent.PollInterval, "i", time.Duration(2)*time.Second, "update interval")
	flag.DurationVar(&cfgAgent.ReportInterval, "i", time.Duration(10)*time.Second, "upload interval")
	flag.StringVar(&cfgAgent.Hash, "k", "", "hash")

	flag.Parse()
	if err := env.Parse(&cfgAgent); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cfgAgent
}