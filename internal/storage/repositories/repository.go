package repositories

import (
	"fmt"

	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

type Repository struct {
	Metrics map[string]storage.Metrics
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewRepository() *Repository {
	metrics := storage.MetricsStore
	gauge := storage.GaugeData
	counter := storage.CounterData
	return &Repository{
		Metrics: metrics,
		Gauge:   gauge,
		Counter: counter,
	}
}

func (r *Repository) CollectMetrics(m storage.Metrics) error {
	r.Metrics[m.ID] = m
	return nil
}

func (r *Repository) GetMetric(metric *storage.Metrics) (storage.Metrics, error) {
	if metric.MType == "gauge" {
		m, ok := r.Metrics[metric.ID]
		if !ok {
			err := fmt.Errorf("id not found %s", metric.ID)
			return *metric, err
		}
		metric.Value = m.Value
		return *metric, nil
	} else if metric.MType == "counter" {
		m, ok := r.Metrics[metric.ID]
		if !ok {
			err := fmt.Errorf("id not found %s", metric.ID)
			return *metric, err
		}
		metric.Delta = m.Delta
		return *metric, nil
	} else {
		err := fmt.Errorf("type %s is wrong", metric.MType)
		return *metric, err
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

//func (r *Repository) GetAll() string {
//	result := "Gauge:\n"
//	for s, f := range r.Gauge {
//		result += fmt.Sprintf("%s : %f\n", s, f)
//	}
//	result = result + "Counter:\n"
//	for s, i := range r.Counter {
//		result += fmt.Sprintf("%s : %d\n", s, i)
//	}
//	return result
//}
