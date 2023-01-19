package storage

// Storage интерфейс хранилища.
type Storage interface {
	GetMetric(metric Metrics) (Metrics, error)               //выгрузка метрики.
	GetAll() (string, error)                                 //выгрузка всех метрик.
	CollectMetric(m Metrics) error                           //получение метрики.
	CollectOrChangeGauge(name string, value float64) error   //получение или изменение метрики типа Gauge.
	CollectOrIncreaseCounter(name string, value int64) error //получение или изменение метрики типа Counter.
	GetGauge(name string) (float64, error)                   //выгрузка метрики типа Gauge.
	GetCounter(name string) (int64, error)                   //выгрузка метрики типа Counter.
	PingClient() error                                       //ping
	CollectMetrics(m []Metrics) error                        //получение метрики батчами.
}
