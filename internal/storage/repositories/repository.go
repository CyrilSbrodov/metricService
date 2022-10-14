package repositories

import (
	"fmt"

	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

type Repository struct {
	Metrics map[string]storage.Metrics
}

func NewRepository() *Repository {
	metrics := storage.MetricsStore
	return &Repository{
		Metrics: metrics,
	}
}

func (r *Repository) CollectMetrics(store map[string]storage.Metrics) error {
	for id, m := range store {
		r.Metrics[id] = m
	}
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
