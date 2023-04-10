# Сервис сбора метрик и алертинга.

Сервис позволяет собирать метрики и отрпавлять их на сервер.

# Начало работы

1. Необходимо запустить сервер.
Предусмотрены различные конфигурации. Флаги, конфигурационный файл.
```GO
type ServerConfig struct {
	Addr             string `json:"address" env:"ADDRESS"` // адрес сервера
	GRPCAddr         string `json:"grpc_addr" env:"GRPC_ADDRESS"` // gRPC адрес сервера
	StoreFile        string `json:"store_file" env:"STORE_FILE"` // файл восстановления значения метрик после перезагрузки сервера.
	Hash             string `env:"KEY"` // хэш
	DatabaseDSN      string `json:"database_dsn" env:"DATABASE_DSN"`// адрес базы данных
	CryptoPROKey     string `json:"crypto_key" env:"CRYPTO_KEY"` // имя крипто файла
	CryptoPROKeyPath string `json:"crypto_pro_key_path" env:"CRYPTO_KEY_PATH"` // путь до крипто файла
	TrustedSubnet    string `json:"trusted_subnet" env:"TRUSTED_SUBNET"` // разрешенный IP
	Config           string // имя файла конфигурации
	Restore          bool          `json:"restore" env:"RESTORE"` // флаг восстановления значения метрик из файла
	StoreInterval    time.Duration `json:"store_interval" env:"STORE_INTERVAL"` // интервал сохранения метрик в файл.
}
```
3. Зпустить агент по сбору метрик.
Предусмотрены различные конфигурации. Флаги, конфигурационный файл.
```GO
type AgentConfig struct {
	Addr             string `json:"address" env:"ADDRESS"`
	GRPCAddr         string `json:"grpc_addr" env:"GRPC_ADDRESS"`
	Config           string
	Hash             string        `env:"KEY"`
	CryptoPROKey     string        `json:"crypto_key" env:"CRYPTO_KEY"`
	CryptoPROKeyPath string        `json:"crypto_key_path" env:"CRYPTO_KEY_PATH"`
	TrustedSubnet    string        `json:"trusted_subnet" env:"TRUSTED_SUBNET"`
	PollInterval     time.Duration `json:"poll_interval" env:"POLL_INTERVAL"`
	ReportInterval   time.Duration `json:"report_interval" env:"REPORT_INTERVAL"`
}
```
