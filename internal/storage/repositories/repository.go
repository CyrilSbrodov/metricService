package repositories

import (
	"fmt"

	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

type Repository struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewRepository() *Repository {
	gauge := storage.GaugeData
	counter := storage.CounterData
	return &Repository{
		Gauge:   gauge,
		Counter: counter,
	}
}

func (r *Repository) CollectGauge(name string, value float64) error {

	r.Gauge[name] = value

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
	value, ok := r.Gauge[name]
	if !ok {
		return value, fmt.Errorf("missing metric %s", name)
	}
	return value, nil
}

func (r *Repository) GetCounter(name string) (int64, error) {
	value, ok := r.Counter[name]
	if !ok {
		return value, fmt.Errorf("missing metric %s", name)
	}
	return value, nil
}

func (r *Repository) GetAll() string {
	result := "Gauge:\n"
	for s, f := range r.Gauge {
		result += fmt.Sprintf("%s : %f\n", s, f)
	}
	result = result + "Counter:\n"
	for s, i := range r.Counter {
		result += fmt.Sprintf("%s : %d\n", s, i)
	}
	return result
}
