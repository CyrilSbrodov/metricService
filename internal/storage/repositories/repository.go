package repositories

import (
	"bufio"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

//создание структуры репозитория
type Repository struct {
	Metrics       map[string]storage.Metrics
	Gauge         map[string]float64
	Counter       map[string]int64
	file          *os.File
	Check         bool
	StoreInterval time.Duration
	Hash          string
	Dsn           string
}

func NewRepository(cfg *config.ServerConfig) (*Repository, error) {
	metrics := storage.MetricsStore
	gauge := storage.GaugeData
	counter := storage.CounterData

	//определение сбора данных из файла
	if cfg.Restore {
		err := restore(&metrics, cfg)
		if err != nil {
			return nil, err
		}
	}

	//определение записи на диск и создание файла
	if cfg.StoreFile == "" {
		return &Repository{
			Metrics:       metrics,
			Gauge:         gauge,
			Counter:       counter,
			file:          nil,
			Check:         false,
			StoreInterval: cfg.StoreInterval,
			Hash:          cfg.Hash,
			Dsn:           cfg.DatabaseDSN,
		}, nil
	} else {
		file, err := os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			return nil, err
		}
		return &Repository{
			Metrics:       metrics,
			Gauge:         gauge,
			Counter:       counter,
			file:          file,
			Check:         true,
			StoreInterval: cfg.StoreInterval,
			Hash:          cfg.Hash,
			Dsn:           cfg.DatabaseDSN,
		}, nil
	}
}

func (r *Repository) CollectMetrics(m storage.Metrics) error {

	if m.Hash != "" {
		_, ok := hashing(r.Hash, &m)
		if !ok {
			err := fmt.Errorf("hash is wrong")
			return err
		}
	}

	switch m.MType {
	case "counter":
		entry, ok := r.Metrics[m.ID]
		if !ok {
			r.Metrics[m.ID] = m
			return nil
		}
		*entry.Delta += *m.Delta
		entry.Hash = m.Hash
		r.Metrics[m.ID] = entry
		return nil
	case "gauge":
		r.Metrics[m.ID] = m
		return nil
	}
	r.Metrics[m.ID] = m
	return nil
}

func (r *Repository) GetMetric(m storage.Metrics) (storage.Metrics, error) {
	if m.MType == "gauge" || m.MType == "counter" {
		metric, ok := r.Metrics[m.ID]
		if !ok {
			err := fmt.Errorf("id not found %s", m.ID)
			return m, err
		}
		metric.Hash, _ = hashing(r.Hash, &metric)
		return metric, nil
	} else {
		err := fmt.Errorf("type %s is wrong", m.MType)
		return m, err
	}
}

//получение всех метрик
func (r *Repository) GetAll() string {
	result := ""
	for s, f := range r.Metrics {
		if f.MType == "gauge" {
			if f.Value != nil {
				result += fmt.Sprintf("%s : %f\n", s, *f.Value)
			}
			continue
		} else if f.MType == "counter" {
			result += fmt.Sprintf("%s : %d\n", s, *f.Delta)
		}
	}
	return result
}

func (r *Repository) CollectOrChangeGauge(name string, value float64) error {
	r.Gauge[name] = value
	var m storage.Metrics
	m.ID = name
	m.Value = &value
	r.Metrics[name] = m

	return nil
}

func (r *Repository) CollectOrIncreaseCounter(name string, value int64) error {
	val, ok := r.Counter[name]
	if !ok {
		r.Counter[name] = value
		return nil
	}
	r.Counter[name] = value + val

	return nil
}

func (r *Repository) GetGauge(name string) (float64, error) {
	value, ok := r.Metrics[name]
	if !ok {
		return 0, fmt.Errorf("missing metric %s", name)
	}
	return *value.Value, nil
}

func (r *Repository) GetCounter(name string) (int64, error) {
	value, ok := r.Counter[name]
	if !ok {
		return value, fmt.Errorf("missing metric %s", name)
	}
	return value, nil
}

func (r *Repository) PingClient(ctx context.Context) error {
	fmt.Println("try to ping DB")
	pool, err := pgx.Connect(ctx, r.Dsn)
	if err != nil {
		fmt.Println("lost connection")
		fmt.Println(err)
		return err
	}
	defer pool.Close(ctx)

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = pool.Ping(ctxTimeout); err != nil {
		fmt.Println("Not ping")
		fmt.Println(err)
		return err
	}
	fmt.Println("start and ping")
	return nil
}

//функция забора данных из файла при запуске
func restore(store *map[string]storage.Metrics, cfg *config.ServerConfig) error {

	file, err := os.OpenFile(cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		err = json.Unmarshal(data, &store)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	defer file.Close()
	return nil
}

//функция загрузки данных на диск
func (r *Repository) Upload() error {
	data, err := json.Marshal(&r.Metrics)
	if err != nil {
		return err
	}
	// записываем событие в буфер
	writer := bufio.NewWriter(r.file)
	if _, err := writer.Write(data); err != nil {
		return err
	}
	if err := writer.WriteByte('\n'); err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func hashing(hashKey string, m *storage.Metrics) (string, bool) {
	var hash string
	switch m.MType {
	case "counter":
		hash = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
	case "gauge":
		hash = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
	}
	h := hmac.New(sha256.New, []byte(hashKey))
	h.Write([]byte(hash))
	hashAccept, err := hex.DecodeString(m.Hash)
	if err != nil {
		return "", false
	}
	return fmt.Sprintf("%x", h.Sum(nil)), hmac.Equal(h.Sum(nil), hashAccept)
}
