#### Оглавление:
____
0. [Сервис сбора метрик и алертинга](https://github.com/CyrilSbrodov/metricService#Сервис-сбора-метрик-и-алертинга).
1. [Начало работы](https://github.com/CyrilSbrodov/metricService#Начало-работы).
1.1. [Конфигурация](https://github.com/CyrilSbrodov/metricService#Конфигурация).
1.2. [Зависимости](https://github.com/CyrilSbrodov/metricService#Зависимости).
2. [Агент](https://github.com/CyrilSbrodov/metricService#Агент).
3. [Сервер](https://github.com/CyrilSbrodov/metricService#Сервер).
____

# Сервис сбора метрик и алертинга.

Сервис позволяет собирать метрики компьютера и отрпавлять их на сервер в зашифрованном виде.
```GO
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики.
	MType string   `json:"type"`            // параметр, принимающий значение Gauge или Counter.
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи Counter.
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи Gauge.
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции.
}
```
Виды [метрик](https://cs.opensource.google/go/go/+/go1.20.3:src/runtime/mstats.go;l=58):
```GO
type MemStats struct {
	Alloc uint64
	TotalAlloc uint64
	Sys uint64
	Lookups uint64
	Mallocs uint64
	//etc.
}
```

# 1. Начало работы

## 1.1. Конфигурация

Предусмотрены различные конфигурации. Флаги, облачные переменные, конфигурационный файл.
Флаги при запуске сервера:
```
-a //адрес сервера
-i //интервал сохранения метрик в файл
-f //файл сохранения метрик
-r //включение восстановления метрик из файла
-k //hash key
-d //адрес бд
-crypto-key //имя файла ключа шифрования
-crypto-key-path //путь файла шифрования
-t //CIDR
-grpc //адрес grpc
-config //конфигурационный файл
```
Флаги при запуске агента:
```
-a //адрес сервера
-p //интервал обновления метрик
-r //интервал отправки метрик на сервер
-k //hash key
-crypto-key //имя файла ключа шифрования
-crypto-key-path //путь файла шифрования
-t //CIDR
-grpc //адрес grpc
```

Склонируйте репозиторий с github:
```
git clone https://github.com/CyrilSbrodov/metricService.git
```
Пример запуска с адресом сервера:
```
cd cmd/server
go run main.go -a localhost:8080
```

Либо создайте конфиг файл в формате json
```
{
    "address": "localhost:8080"
    //etc.
}
```

1. Необходимо запустить сервер из пакета [cmd](https://github.com/CyrilSbrodov/metricService/blob/main/cmd/server/main.go)
```
go run main.go
```

2. Зпустить агент из пакета [cmd](https://github.com/CyrilSbrodov/metricService/blob/main/cmd/agent/main.go)
```
go run main.go
```
## 1.2. Зависимости.
Для работы обязательно понадобится PostgreSQL последней версии.


# 2. Агент

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

# 3. Сервер
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
