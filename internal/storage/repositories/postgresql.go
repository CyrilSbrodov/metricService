package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"github.com/CyrilSbrodov/metricService.git/pkg/client/postgresql"
)

type PGSStore struct {
	client postgresql.Client
}

func createTable(ctx context.Context, client postgresql.Client) error {
	tx, err := client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	q := `CREATE TABLE metrics (
    id TEXT NOT NULL,
    mType TEXT NOT NULL,
    delta INT,
    value DOUBLE PRECISION,
    hash TEXT
);`
	_, err = tx.Exec(ctx, q)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func NewPGSStore(client postgresql.Client) (*PGSStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := createTable(ctx, client); err != nil {
		return nil, err
	}
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
