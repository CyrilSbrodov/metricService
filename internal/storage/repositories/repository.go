package repositories

import (
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

	//fmt.Println("r.Gauge")
	//fmt.Println(r.Gauge)
	return nil
}

func (r *Repository) CollectCounter(name string, value int64) error {

	r.Counter[name] = value

	//fmt.Println("r.Counter")
	//fmt.Println(r.Counter)
	return nil
}
