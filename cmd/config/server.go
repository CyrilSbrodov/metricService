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
	Addr             string `json:"address" env:"ADDRESS"`
	GRPCAddr         string `json:"grpc_addr" env:"GRPC_ADDRESS"`
	Config           string
	StoreFile        string        `json:"store_file" env:"STORE_FILE"`
	Hash             string        `env:"KEY"`
	DatabaseDSN      string        `json:"database_dsn" env:"DATABASE_DSN"`
	CryptoPROKey     string        `json:"crypto_key" env:"CRYPTO_KEY"`
	CryptoPROKeyPath string        `json:"crypto_pro_key_path" env:"CRYPTO_KEY_PATH"`
	TrustedSubnet    string        `json:"trusted_subnet" env:"TRUSTED_SUBNET"`
	Restore          bool          `json:"restore" env:"RESTORE"`
	StoreInterval    time.Duration `json:"store_interval" env:"STORE_INTERVAL"`
}

// ServerConfigInit инициализация конфига.
func ServerConfigInit() *ServerConfig {
	cfgSrv := &ServerConfig{}
	path := getPathServerConfigFile()
	cfgSrv.Config = path
	if cfgSrv.Config != "" {
		configFile, err := os.Open(cfgSrv.Config)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer configFile.Close()
		jsonParser := json.NewDecoder(configFile)
		jsonParser.Decode(&cfgSrv)
	}
	flag.StringVar(&cfgSrv.Addr, "a", "localhost:8080", "ADDRESS")
	flag.DurationVar(&cfgSrv.StoreInterval, "i", time.Duration(300)*time.Second, "STORE_INTERVAL")
	flag.StringVar(&cfgSrv.StoreFile, "f", "/tmp/devops-metrics-db.json", "STORE_FILE")
	flag.BoolVar(&cfgSrv.Restore, "r", true, "RESTORE")
	flag.StringVar(&cfgSrv.Hash, "k", "", "KEY")
	flag.StringVar(&cfgSrv.DatabaseDSN, "d", "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable", "DATABASE_DSN")
	flag.StringVar(&cfgSrv.CryptoPROKey, "crypto-key", "", "path to file")
	flag.StringVar(&cfgSrv.CryptoPROKeyPath, "crypto-key-path", "../../cmd/server/", "path to folder")
	flag.StringVar(&cfgSrv.TrustedSubnet, "t", "", "CIDR")
	flag.StringVar(&cfgSrv.GRPCAddr, "grpc", ":3200", "grpc port")
	flag.Parse()
	if err := env.Parse(cfgSrv); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cfgSrv
}

//getPathServerConfigFile загрузка конфига из файла
func getPathServerConfigFile() string {
	var path string
	flag.StringVar(&path, "c/-config", "", "path to config file")
	path = os.Getenv("CONFIG")

	return path
}
