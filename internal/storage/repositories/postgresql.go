package repositories

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
)

type DB struct {
	databaseURL string
	db          *pgx.Conn
}

func NewDB(cfg *config.ServerConfig) (*DB, error) {
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	return &DB{
		databaseURL: cfg.DatabaseDSN,
		db:          conn,
	}, nil
}

func (db *DB) Connect(ctx context.Context) error {
	if err := db.db.Ping(ctx); err != nil {
		panic(err)
	}
	return nil
}

//func (db *DB) Ping() {
//	conn.q
//}
