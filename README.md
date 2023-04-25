#### Оглавление:
____
0. [Сервис сбора метрик и алертинга](https://github.com/CyrilSbrodov/metricService#Сервис-сбора-метрик-и-алертинга).
1. [ЗАВИСИМОСТИ](https://github.com/CyrilSbrodov/metricService#ЗАВИСИМОСТИ).
2. [ЗАПУСК/СБОРКА](#запусксборка).
2.1. [Конфигурация](https://github.com/CyrilSbrodov/metricService#Конфигурация).
2.1.1 [Флаги](https://github.com/CyrilSbrodov/metricService#1-флаги).
2.1.2 [Облачные переменные](https://github.com/CyrilSbrodov/metricService#2-облачные-переменные).
2.1.3 [Конфигурационный файл](https://github.com/CyrilSbrodov/metricService#3-конфигурационный-файл).
2.2. [Запуск сервера](https://github.com/CyrilSbrodov/metricService#Запуск-сервера).
2.3. [Запуск агента](https://github.com/CyrilSbrodov/metricService#Запуск-агента).
3. [Для разработчиков](https://github.com/CyrilSbrodov/metricService#Для-разработчиков).
3.1. [Агент](https://github.com/CyrilSbrodov/metricService#1-агент).
3.2. [Сервер](https://github.com/CyrilSbrodov/metricService#2-Сервер).
____

# Сервис сбора метрик и алертинга.

Сервис позволяет собирать метрики ПК в системном формате (CPU, RAM, HDD)(числовые метрики) и передает их по протоколам gRPC и HTML зашифрованном формате в БД PostgreSQL.
Возможность использовать свои метрики.

Структура [метрик](https://github.com/CyrilSbrodov/metricService/blob/main/internal/storage/models.go):
```GO
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики.
	MType string   `json:"type"`            // параметр, принимающий значение Gauge или Counter.
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи Counter.
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи Gauge.
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции.
}
```
Виды собираемых [метрик](https://cs.opensource.google/go/go/+/go1.20.3:src/runtime/mstats.go;l=58):
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

Структура сервиса следующая:
1) Сервер - обработка полученных данных и отправка их в БД Postgres.
2) Агент - сбор данных с ПК.
3) БД - прием получаемых данных.
____
# ЗАВИСИМОСТИ.

Используется язык go версии 1.18. Используемые библиотеки:
- github.com/caarlos0/env/v6 v6.10.1
- github.com/go-chi/chi/v5 v5.0.7
- github.com/jackc/pgx/v5 v5.0.4
- github.com/lib/pq v1.10.7
- github.com/rs/zerolog v1.28.0
- github.com/shirou/gopsutil/v3 v3.22.10
- github.com/stretchr/testify v1.8.1
- golang.org/x/tools v0.5.0
- honnef.co/go/tools v0.3.3
- POSTGRESQL v17
____

# ЗАПУСК/СБОРКА

## Конфигурация

Предусмотрены различные конфигурации:
1) флаги
2) облачные переменные 
3) конфигурационный файл

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

## 1) Флаги
Параметры запуска передаются в формате: -a localhost:8080, 
где "-а" параметр адреса сервера, "localhost:8080" адрес сервера.

Это позволяет запускать утилиту следующим образом:
```
cd cmd/server
go run main.go -a localhost:8080
```

## 2) Облачные переменные
Перед запуском утилиты необходимо присвоить переменным значения в формате: 
```
ADDRESS='localhost:8080'
```
где "ADDRESS" - облачная переменная, 'localhost:8080' присвоение значения данной переменной.

## 3) Конфигурационный файл
Для исрользования конфиг файла нужен следующий формат:
```
{
    "address": "localhost:8080"
    //etc.
}
```
## Запуск сервера

Необходимо запустить сервер из пакета [cmd](https://github.com/CyrilSbrodov/metricService/blob/main/cmd/server/main.go)
```
cd cmd/server
go run main.go
```

## Запуск агента

Необходимо запустить агент из пакета [cmd](https://github.com/CyrilSbrodov/metricService/blob/main/cmd/agent/main.go)
```
cd cmd/agent
go run main.go
```
____
 
# Для разработчиков
Структура приложения позволяет нативно вносить корректировки:
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

Немного о сервере и агенте:

# 1. Агент

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

# 2. Сервер
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
