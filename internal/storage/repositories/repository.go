package repositories

import (
	"fmt"
	"strconv"
)

type Repository struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewRepository() *Repository {
	gauge := make(map[string]float64)
	counter := make(map[string]int64)
	return &Repository{
		Gauge:   gauge,
		Counter: counter,
	}
}

func (r *Repository) Collect(types, name, value string) error {
	if types == "gauge" {
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err //TODO
		}
		r.Gauge[name] = value
	} else if types == "counter" {
		value, err := strconv.Atoi(value)
		if err != nil {
			return err //TODO
		}
		r.Counter[name] = int64(value)
	}
	fmt.Println("r.Gauge")
	fmt.Println(r.Gauge)
	fmt.Println("r.Counter")
	fmt.Println(r.Counter)
	return nil
}
