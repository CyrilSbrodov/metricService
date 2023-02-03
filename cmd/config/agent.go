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
	Addr           string `json:"address" env:"ADDRESS"`
	Config         string
	Hash           string        `env:"KEY"`
	CryptoPROKey   string        `json:"crypto_key" env:"CRYPTO_KEY"`
	PollInterval   time.Duration `json:"poll_interval" env:"POLL_INTERVAL"`
	ReportInterval time.Duration `json:"report_interval" env:"REPORT_INTERVAL"`
}

// AgentConfigInit инициализая конфига.
func AgentConfigInit() *AgentConfig {
	cfgAgent := &AgentConfig{}
	path := getPathAgentConfigFile()
	cfgAgent.Config = path
	if cfgAgent.Config != "" {
		configFile, err := os.Open(cfgAgent.Config)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer configFile.Close()
		jsonParser := json.NewDecoder(configFile)
		jsonParser.Decode(&cfgAgent)
	}
	flag.StringVar(&cfgAgent.Addr, "a", "localhost:8080", "ADDRESS")
	flag.DurationVar(&cfgAgent.PollInterval, "p", time.Duration(2)*time.Second, "update interval")
	flag.DurationVar(&cfgAgent.ReportInterval, "r", time.Duration(10)*time.Second, "upload interval")
	flag.StringVar(&cfgAgent.Hash, "k", "", "hash")
	flag.StringVar(&cfgAgent.CryptoPROKey, "crypto-key", "", "path to file")
	flag.Parse()
	if err := env.Parse(cfgAgent); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cfgAgent
}

//getPathAgentConfigFile загрузка конфига из файла
func getPathAgentConfigFile() string {
	var path string
	flag.StringVar(&path, "c/-config", "", "path to config file")
	path = os.Getenv("CONFIG")

	return path
}
