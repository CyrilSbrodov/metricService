package storage

type Storage interface {
	GetMetric(metric *Metrics) (Metrics, error)
	GetAll() string
	CollectMetrics(store map[string]Metrics) error
}
