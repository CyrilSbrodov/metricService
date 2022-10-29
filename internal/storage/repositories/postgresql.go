package repositories

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
)

type DB struct {
	databaseURL string
	//db          *pgxpool.Pool
}

func NewDB(cfg *config.ServerConfig) (*DB, error) {
	//var pool *pgxpool.Pool

	return &DB{
		databaseURL: cfg.DatabaseDSN,
		//db:          pool,
	}, nil
}

func (db *DB) Connect() error {
	pool, err := sql.Open("postgres", db.databaseURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(); err != nil {
		return err
	}
	return nil
}

//func (db *DB) Ping() {
//	conn.q
//}
