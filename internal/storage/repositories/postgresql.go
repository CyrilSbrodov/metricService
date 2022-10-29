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
	dbpgx       *pgx.Conn
}

func NewDB(cfg *config.ServerConfig) (*DB, error) {
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
		return nil, err
	}
	defer conn.Close(context.Background())
	return &DB{
		databaseURL: cfg.DatabaseDSN,
		dbpgx:       conn,
	}, nil
}

func (db *DB) Connect() error {
	err := db.dbpgx.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}

//func (db *DB) Ping() {
//	conn.q
//}
