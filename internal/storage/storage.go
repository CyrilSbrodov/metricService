package storage

type Storage interface {
	Collect(types, name, value string) error
}
