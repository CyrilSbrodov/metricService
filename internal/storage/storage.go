package storage

type Storage interface {
	GetMetric(metric Metrics) (Metrics, error)
	GetAll() string
	CollectMetrics(m Metrics) error
	CollectOrChangeGauge(name string, value float64) error
	CollectOrIncreaseCounter(name string, value int64) error
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
}
