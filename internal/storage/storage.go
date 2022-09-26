package storage

type Storage interface {
	CollectGauge(name, value string) error
	CollectCounter(name, value string) error
}
