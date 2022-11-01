package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/cmd/server/client/postgresql"
)

//type db struct {
//	client postgresql.Client
//}

//func NewDB(client postgresql.Client) storage.Storage {
//	return &db{
//		client: client,
//	}
//}

//func (d db) GetMetric(m storage.Metrics) (storage.Metrics, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (d db) GetAll() string {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (d db) CollectMetrics(m storage.Metrics) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (d db) CollectOrChangeGauge(name string, value float64) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (d db) CollectOrIncreaseCounter(name string, value int64) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (d db) GetGauge(name string) (float64, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (d db) GetCounter(name string) (int64, error) {
//	//TODO implement me
//	panic("implement me")
//}

type DB struct {
	databaseURL string
	client      postgresql.Client
}

func NewDB(cfg *config.ServerConfig, client postgresql.Client) (*DB, error) {
	//var pool *pgxpool.Pool

	return &DB{
		databaseURL: cfg.DatabaseDSN,
		client:      client,
	}, nil
}

func (db *DB) PingClient() error {
	if err := db.client.Ping(context.Background()); err != nil {
		return err
	}
	//pool, err := sql.Open("postgres", db.databaseURL)
	//
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	//	os.Exit(1)
	//}
	//defer pool.Close()
	//
	//if err := pool.Ping(); err != nil {
	//	return err
	//}
	return nil
}

//func (db *DB) Ping() {
//	conn.q
//}

func (db *DB) Ping() error {
	pool, err := sql.Open("postgres", db.databaseURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err = pool.Ping(); err != nil {
		return err
	}
	return nil
}
