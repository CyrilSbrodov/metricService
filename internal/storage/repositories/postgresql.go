package repositories

import (
	"context"

	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"github.com/CyrilSbrodov/metricService.git/pkg/client/postgresql"
)

type PGSStore struct {
	client postgresql.Client
}

func NewPGSStore(client postgresql.Client) (*PGSStore, error) {
	return &PGSStore{
		client: client,
	}, nil
}

func (p *PGSStore) GetMetric(m storage.Metrics) (storage.Metrics, error) {
	//TODO implement me

	return storage.Metrics{}, nil
}

func (p *PGSStore) GetAll() string {
	//TODO implement me
	return ""
}

func (p *PGSStore) CollectMetrics(m storage.Metrics) error {
	//TODO implement me
	return nil
}

func (p *PGSStore) CollectOrChangeGauge(name string, value float64) error {
	//TODO implement me
	return nil
}

func (p *PGSStore) CollectOrIncreaseCounter(name string, value int64) error {
	//TODO implement me
	return nil
}

func (p *PGSStore) GetGauge(name string) (float64, error) {
	//TODO implement me
	return 0, nil
}

func (p *PGSStore) GetCounter(name string) (int64, error) {
	//TODO implement me
	return 0, nil
}

func (p *PGSStore) PingClient() error {
	return p.client.Ping(context.Background())
}
