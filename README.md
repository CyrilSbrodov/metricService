#### Оглавление:
____
0. [Сервис сбора метрик и алертинга](https://github.com/CyrilSbrodov/metricService#Сервис-сбора-метрик-и-алертинга).
1. [Начало работы](https://github.com/CyrilSbrodov/metricService#Начало-работы).
2. [Агент](https://github.com/CyrilSbrodov/metricService#Агент).
3. [Сервер](https://github.com/CyrilSbrodov/metricService#Сервер).
____

# Сервис сбора метрик и алертинга.

Сервис позволяет собирать метрики и отрпавлять их на сервер.

# Начало работы

1. Необходимо запустить сервер.

Предусмотрены различные конфигурации. Флаги, конфигурационный файл.

[ServerConfig](https://github.com/CyrilSbrodov/metricService/blob/main/cmd/config/server.go):
```GO
type ServerConfig struct {
	Addr             string `json:"address" env:"ADDRESS"` // адрес сервера
	GRPCAddr         string `json:"grpc_addr" env:"GRPC_ADDRESS"` // gRPC адрес сервера
	StoreFile        string `json:"store_file" env:"STORE_FILE"` // файл восстановления значения метрик после перезагрузки сервера.
	Hash             string `env:"KEY"` // хэш ключ
	DatabaseDSN      string `json:"database_dsn" env:"DATABASE_DSN"`// адрес базы данных
	CryptoPROKey     string `json:"crypto_key" env:"CRYPTO_KEY"` // имя крипто файла
	CryptoPROKeyPath string `json:"crypto_pro_key_path" env:"CRYPTO_KEY_PATH"` // путь до крипто файла
	TrustedSubnet    string `json:"trusted_subnet" env:"TRUSTED_SUBNET"` // разрешенный IP
	Config           string // имя файла конфигурации
	Restore          bool          `json:"restore" env:"RESTORE"` // флаг восстановления значения метрик из файла
	StoreInterval    time.Duration `json:"store_interval" env:"STORE_INTERVAL"` // интервал сохранения метрик в файл.
}
```
2. Зпустить агент по сбору метрик.
 
Предусмотрены различные конфигурации. Флаги, конфигурационный файл.

[AgentConfig](https://github.com/CyrilSbrodov/metricService/blob/main/cmd/config/agent.go):
```GO
type AgentConfig struct {
	Addr             string `json:"address" env:"ADDRESS"` // адрес сервера
	GRPCAddr         string `json:"grpc_addr" env:"GRPC_ADDRESS"` // gRPC адрес
	Config           string // имя файла конфигурации
	Hash             string        `env:"KEY"` // хэш ключ
	CryptoPROKey     string        `json:"crypto_key" env:"CRYPTO_KEY"` // имя крипто файла
	CryptoPROKeyPath string        `json:"crypto_key_path" env:"CRYPTO_KEY_PATH"` // путь до крипто файла
	TrustedSubnet    string        `json:"trusted_subnet" env:"TRUSTED_SUBNET"` // разрешенный IP
	PollInterval     time.Duration `json:"poll_interval" env:"POLL_INTERVAL"` // интервал обновления метрик
	ReportInterval   time.Duration `json:"report_interval" env:"REPORT_INTERVAL"` // интервал отправки метрик на сервер
}
```

# Агент

[Структура агента](https://github.com/CyrilSbrodov/metricService/blob/main/internal/app/agent.go):
```GO
type AgentApp struct {
	client *http.Client // клиент
	cfg    config.AgentConfig // конфиг
	logger *loggers.Logger // логгер
	public *rsa.PublicKey // публичный ключ шифрования
	url    string 
	wg     sync.WaitGroup
}
```
Метрики собираются с интерваорм, согласно конфигу (по умолчанию интервал составляет 2 секунды).
Метрики отправляются на сервер по одной и батчами с интервалом, согласно конфигу (по умолчанию интервал составляет 10 секунд).

Возможность выбора отправки метрик по протоколам http или gRPC (если в конфиге указан адрес gRPC, то метрики автоматически отправляются по этому протоколу (по умаолчанию отправка по http))
```GO
flag.StringVar(&cfgAgent.GRPCAddr, "grpc", "", "grpc port")
```

Возможность выбора шифрования (если в конфиге указано имя крипто файла, то все метрики шифруются перед отправкой на сервер (по умолчанию шифрование отключено)).
```GO
flag.StringVar(&cfgAgent.CryptoPROKey, "crypto-key", "", "crypto file")
```

# Сервер
[Структура сервера](https://github.com/CyrilSbrodov/metricService/blob/main/internal/app/server.go):
```GO
type ServerApp struct {
	router   *chi.Mux // роутер
	cfg      config.ServerConfig // конфиг
	logger   *loggers.Logger // логгер
	cryptoer crypto.Cryptoer // интерфейс шифрования
	private  *rsa.PrivateKey // приватный ключ шифрования
}
```
Сервер получает метрики по следующим эндпоинтам:
1) [http](https://github.com/CyrilSbrodov/metricService/blob/main/internal/handlers/handler.go):
```GO
func (h *Handler) Register(r *chi.Mux) {
	r.Post("/value/", trustedSubnet(h.cfg, gzipHandle(h.GetHandlerJSON())))
	r.Get("/value/*", trustedSubnet(h.cfg, gzipHandle(h.GetHandler())))
	r.Get("/", gzipHandle(h.GetAllHandler()))
	r.Post("/update/", trustedSubnet(h.cfg, gzipHandle(h.CollectHandler())))
	r.Post("/update/gauge/*", trustedSubnet(h.cfg, gzipHandle(h.GaugeHandler())))
	r.Post("/update/counter/*", trustedSubnet(h.cfg, gzipHandle(h.CounterHandler())))
	r.Post("/*", trustedSubnet(h.cfg, gzipHandle(h.OtherHandler())))
	r.Get("/ping", h.PingDB())
	r.Post("/updates/", trustedSubnet(h.cfg, gzipHandle(h.CollectBatchHandler())))
	r.Mount("/debug", middleware.Profiler())
}
```
2) [gRPC](https://github.com/CyrilSbrodov/metricService/blob/main/internal/app/protoserver/server.go):
```GO
//получение метрик по одной
func (s *StoreServer) CollectMetric(ctx context.Context, in *pb.AddMetricRequest) (*pb.AddMetricResponse, error) {
	var response pb.AddMetricResponse
	var m storage.Metrics
	metric, err := s.Storage.GetMetric(*m.ToMetric(in.Metrics))
	if err != nil {
		s.logger.LogErr(err, "")
		return nil, err
	}
	response.Metrics = metric.ToProto()
	return &response, nil
}

// получение метрик батчами
func (s *StoreServer) CollectMetrics(ctx context.Context, in *pb.AddMetricsRequest) (*pb.AddMetricsResponse, error) {
	var response pb.AddMetricsResponse
	var metrics []storage.Metrics
	var m storage.Metrics
	for _, metric := range in.Metrics {
		metrics = append(metrics, *m.ToMetric(metric))
	}
	err := s.Storage.CollectMetrics(metrics)
	if err != nil {
		s.logger.LogErr(err, "")
		return nil, err
	}
	return &response, nil
}
```
