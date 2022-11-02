package repositories

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

type PGSStore struct {
	db *sql.DB
}

func NewPGSStore(cfg *config.ServerConfig) (storage.Storage, error) {
	db, err := sql.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		fmt.Println("lost connection")
		fmt.Println(err)
		return nil, err
	}
	return &PGSStore{
		db: db,
	}, nil
}

func (P PGSStore) GetMetric(m storage.Metrics) (storage.Metrics, error) {
	//TODO implement me
	panic("implement me")
}

func (P PGSStore) GetAll() string {
	//TODO implement me
	panic("implement me")
}

func (P PGSStore) CollectMetrics(m storage.Metrics) error {
	//TODO implement me
	panic("implement me")
}

func (P PGSStore) CollectOrChangeGauge(name string, value float64) error {
	//TODO implement me
	panic("implement me")
}

func (P PGSStore) CollectOrIncreaseCounter(name string, value int64) error {
	//TODO implement me
	panic("implement me")
}

func (P PGSStore) GetGauge(name string) (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (P PGSStore) GetCounter(name string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (P PGSStore) PingClient(ctx context.Context) error {
	return P.db.Ping()
}
