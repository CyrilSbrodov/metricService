package repositories

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

type Repository struct {
	Metrics       map[string]storage.Metrics
	Gauge         map[string]float64
	Counter       map[string]int64
	file          *os.File
	check         bool
	StoreInterval time.Duration
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	metrics := storage.MetricsStore
	gauge := storage.GaugeData
	counter := storage.CounterData
	if cfg.Restore {
		err := Restore(&metrics, cfg)
		if err != nil {
			return nil, err
		}
	}

	if cfg.StoreFile == "" {
		return &Repository{
			Metrics:       metrics,
			Gauge:         gauge,
			Counter:       counter,
			file:          nil,
			check:         false,
			StoreInterval: cfg.StoreInterval,
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
			check:         true,
			StoreInterval: cfg.StoreInterval,
		}, nil
	}
}

func (r *Repository) CollectMetrics(m storage.Metrics) error {

	if m.MType == "counter" && r.Metrics[m.ID].Delta != nil {
		var val int64
		val = *r.Metrics[m.ID].Delta
		val += *m.Delta
		*r.Metrics[m.ID].Delta = val
		return nil
	} else {
		r.Metrics[m.ID] = m
		return nil
	}
}

func (r *Repository) GetMetric(metric storage.Metrics) (storage.Metrics, error) {
	if metric.MType == "gauge" || metric.MType == "counter" {
		m, ok := r.Metrics[metric.ID]
		if !ok {
			err := fmt.Errorf("id not found %s", metric.ID)
			return metric, err
		}
		return m, nil
	} else {
		err := fmt.Errorf("type %s is wrong", metric.MType)
		return metric, err
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

func Restore(store *map[string]storage.Metrics, cfg *config.Config) error {

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
