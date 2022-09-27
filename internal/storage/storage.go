package storage

type Storage interface {
	CollectGauge(name string, value float64) error
	CollectCounter(name string, value int64) error
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	GetAll() string
}
