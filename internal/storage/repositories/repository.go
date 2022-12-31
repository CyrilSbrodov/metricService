package repositories

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/cmd/loggers"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

// Repository структура репозитория.
type Repository struct {
	Metrics       map[string]storage.Metrics
	Gauge         map[string]float64
	Counter       map[string]int64
	file          *os.File
	Check         bool
	StoreInterval time.Duration
	Hash          string
	Dsn           string
	logger        loggers.Logger
	sync          sync.Mutex
}

// NewRepository создание нового репозитория.
func NewRepository(cfg *config.ServerConfig, logger *loggers.Logger) (*Repository, error) {
	metrics := storage.MetricsStore
	gauge := storage.GaugeData
	counter := storage.CounterData
	check := false

	//определение сбора данных из файла.
	if cfg.Restore {
		err := restore(&metrics, cfg, logger)
		if err != nil {
			logger.LogErr(err, "failed to restore data from file")
			return nil, err
		}
	}

	file, err := newStoreFile(cfg.StoreFile, logger)
	if err != nil {
		logger.LogErr(err, "failed to create file")
		return nil, err
	}

	if cfg.StoreInterval == 0 {
		check = true
	}
	repo := &Repository{
		Metrics:       metrics,
		Gauge:         gauge,
		Counter:       counter,
		file:          file,
		Check:         check,
		StoreInterval: cfg.StoreInterval,
	}

	if file != nil {
		ticker := time.NewTicker(cfg.StoreInterval)
		go uploadWithTicker(ticker, repo, logger)
	}

	return &Repository{
		Metrics:       metrics,
		Gauge:         gauge,
		Counter:       counter,
		file:          file,
		Check:         check,
		StoreInterval: cfg.StoreInterval,
		Hash:          cfg.Hash,
		Dsn:           cfg.DatabaseDSN,
	}, nil
}

// CollectMetric сохранение метрики.
func (r *Repository) CollectMetric(m storage.Metrics) error {
	if m.Hash != "" {
		_, ok := hashing(r.Hash, &m, &r.logger)
		if !ok {
			err := fmt.Errorf("hash is wrong")
			r.logger.LogErr(fmt.Errorf("hash is wrong"), "wrong hash")
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

// CollectMetrics сохранение метрик батчами.
func (r *Repository) CollectMetrics(metrics []storage.Metrics) error {
	for _, m := range metrics {
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
	}
	return nil
}

// GetMetric выгрузка метрики.
func (r *Repository) GetMetric(m storage.Metrics) (storage.Metrics, error) {
	if m.MType == "gauge" || m.MType == "counter" {
		metric, ok := r.Metrics[m.ID]
		if !ok {
			r.logger.LogErr(fmt.Errorf("id not found %s", m.ID), "id not found")
			err := fmt.Errorf("id not found %s", m.ID)
			return m, err
		}
		metric.Hash, _ = hashing(r.Hash, &metric, &r.logger)
		return metric, nil
	} else {
		r.logger.LogErr(fmt.Errorf("type %s is wrong", m.MType), "type %s is wrong")
		err := fmt.Errorf("type %s is wrong", m.MType)
		return m, err
	}
}

// GetAll получение всех метрик.
func (r *Repository) GetAll() (string, error) {
	result := ""
	keys := make([]string, 0, len(r.Metrics))
	for k := range r.Metrics {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if r.Metrics[key].MType == "gauge" {
			result += fmt.Sprintf("%s : %f<br>\n", key, *r.Metrics[key].Value)
		} else if r.Metrics[key].MType == "counter" {
			result += fmt.Sprintf("%s : %d<br>\n", key, *r.Metrics[key].Delta)
		}
	}
	return result, nil
}

// CollectOrChangeGauge сохранение или изменение метрики типа Gauge.
func (r *Repository) CollectOrChangeGauge(name string, value float64) error {
	r.Gauge[name] = value
	var m storage.Metrics
	m.ID = name
	m.Value = &value
	r.Metrics[name] = m

	return nil
}

// CollectOrIncreaseCounter сохранение или изменение метрики типа Counter.
func (r *Repository) CollectOrIncreaseCounter(name string, value int64) error {
	val, ok := r.Counter[name]
	if !ok {
		r.Counter[name] = value
		return nil
	}
	r.Counter[name] = value + val

	return nil
}

// GetGauge выгрузка метрики типа Gauge.
func (r *Repository) GetGauge(name string) (float64, error) {
	value, ok := r.Metrics[name]
	if !ok {
		r.logger.LogErr(fmt.Errorf("missing metric %s", name), "metric not found")
		return 0, fmt.Errorf("missing metric %s", name)
	}
	return *value.Value, nil
}

// GetCounter выгрузка метрики типа Counter.
func (r *Repository) GetCounter(name string) (int64, error) {
	value, ok := r.Counter[name]
	if !ok {
		r.logger.LogErr(fmt.Errorf("missing metric %s", name), "metric not found")
		return value, fmt.Errorf("missing metric %s", name)
	}
	return value, nil
}

// PingClient проверка клиента.
func (r *Repository) PingClient() error {
	db, err := sql.Open("postgres", r.Dsn)
	if err != nil {
		r.logger.LogErr(err, "not connection")
		return err
	}

	return db.Ping()
}

//функция забора данных из файла при запуске.
func restore(store *map[string]storage.Metrics, cfg *config.ServerConfig, logger *loggers.Logger) error {

	file, err := os.OpenFile(cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		logger.LogErr(err, "failed to open file")
		return err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		err = json.Unmarshal(data, &store)
		if err != nil {
			logger.LogErr(err, "failed to unmarshal data")
			return err
		}
	}

	defer file.Close()
	return nil
}

// Upload функция загрузки данных на диск.
func (r *Repository) Upload() error {
	data, err := json.Marshal(&r.Metrics)
	if err != nil {
		r.logger.LogErr(err, "failed to marshal data")
		return err
	}
	// записываем событие в буфер
	writer := bufio.NewWriter(r.file)
	if _, err := writer.Write(data); err != nil {
		r.logger.LogErr(err, "failed to write buffer")
		return err
	}
	if err := writer.WriteByte('\n'); err != nil {
		r.logger.LogErr(err, "failed to write bytes")
		return err
	}
	writer.Flush()
	return nil
}

//функция хеширования.
func hashing(hashKey string, m *storage.Metrics, logger *loggers.Logger) (string, bool) {
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
		logger.LogErr(err, "failed to decode data")
		return "", false
	}
	return fmt.Sprintf("%x", h.Sum(nil)), hmac.Equal(h.Sum(nil), hashAccept)
}

//создание нового файла.
func newStoreFile(filename string, logger *loggers.Logger) (*os.File, error) {
	if len(filename) == 0 {
		return nil, nil
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		logger.LogErr(err, "failed to open/create file")
		return nil, err
	}
	return file, nil
}

//функция загрузки данных в файл по тикеру.
func uploadWithTicker(ticker *time.Ticker, repo *Repository, logger *loggers.Logger) {
	for range ticker.C {
		if err := repo.Upload(); err != nil {
			logger.LogErr(err, "failed to upload metrics to file")
			return
		}
	}
}
