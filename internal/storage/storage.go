package storage

type Storage interface {
	CollectGauge(name string, value float64) error
	CollectCounter(name string, value int64) error
}
