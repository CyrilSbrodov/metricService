package repositories

import (
	"fmt"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"strconv"
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

func (r *Repository) CollectGauge(name, value string) error {

	//_, ok := r.Gauge[name]
	//if !ok {
	//	return fmt.Errorf("%s does not exists ", name)
	//}
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err //TODO
	}
	r.Gauge[name] = val

	fmt.Println("r.Gauge")
	fmt.Println(r.Gauge)
	return nil
}

func (r *Repository) CollectCounter(name, value string) error {

	//_, ok := r.Counter[name]
	//if !ok {
	//	return fmt.Errorf("%s does not exists ", name)
	//}
	val, err := strconv.Atoi(value)
	if err != nil {
		return err //TODO
	}
	r.Counter[name] = int64(val)

	fmt.Println("r.Counter")
	fmt.Println(r.Counter)
	return nil
}
