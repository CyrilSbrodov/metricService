/*
	Package config создание конфиг файла для агента
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

// AgentConfig структура конфига для агента.
type AgentConfig struct {
	Addr           string        `json:"address" env:"ADDRESS"`
	Config         string        `env:"CONFIG"`
	Hash           string        `env:"KEY"`
	CryptoPROKey   string        `json:"crypto_key" env:"CRYPTO_KEY"`
	PollInterval   time.Duration `json:"poll_interval" env:"POLL_INTERVAL"`
	ReportInterval time.Duration `json:"report_interval" env:"REPORT_INTERVAL"`
}

var cfgAgent AgentConfig

// AgentConfigInit инициализая конфига.
func AgentConfigInit() AgentConfig {
	cfg := parseFromAgentConfigFile()
	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "ADDRESS")
	flag.DurationVar(&cfg.PollInterval, "p", time.Duration(2)*time.Second, "update interval")
	flag.DurationVar(&cfg.ReportInterval, "r", time.Duration(10)*time.Second, "upload interval")
	flag.StringVar(&cfg.Hash, "k", "", "hash")
	flag.StringVar(&cfg.CryptoPROKey, "crypto-key", "", "path to file")
	//public.pem
	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cfg
}

//parseFromAgentConfigFile загрузка конфига из файла
func parseFromAgentConfigFile() AgentConfig {
	flag.StringVar(&cfgAgent.Config, "c/-config", "", "path to config file")
	flag.Parse()
	if err := env.Parse(&cfgAgent); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if cfgAgent.Config != "" {
		configFile, err := os.Open(cfgAgent.Config)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer configFile.Close()
		jsonParser := json.NewDecoder(configFile)
		jsonParser.Decode(&cfgAgent)
		return cfgAgent
	}
	return cfgAgent
}
