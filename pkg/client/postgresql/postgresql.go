package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/cmd/loggers"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Ping(ctx context.Context) error
}

func NewClient(ctx context.Context, maxAttempts int, cfg *config.ServerConfig, logger *loggers.Logger) (pool *pgxpool.Pool, err error) {
	err = DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

		defer cancel()

		pool, err = pgxpool.New(ctx, cfg.DatabaseDSN)
		if err != nil {
			logger.LogErr(err, "Failure to connect to PostgreSQL")
		}
		return nil

	}, maxAttempts, 5*time.Second)
	if err != nil {
		logger.LogErr(err, "Failure to connect to PostgreSQL")
	}
	return
}
func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}
		return nil
	}
	return
}
