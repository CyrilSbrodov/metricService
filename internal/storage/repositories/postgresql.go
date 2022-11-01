package repositories

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/cmd/server/client/postgresql"
)

type DB struct {
	databaseURL string
	client      postgresql.Client
}

func NewDB(cfg *config.ServerConfig, client postgresql.Client) (*DB, error) {
	return &DB{
		databaseURL: cfg.DatabaseDSN,
		client:      client,
	}, nil
}

func (db *DB) PingClient() error {
	fmt.Println("try to ping DB")
	err := db.client.Ping(context.Background())
	if err != nil {
		fmt.Println("lost connection")
		fmt.Println(err)
		return err
	}
	return nil
}
