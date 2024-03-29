package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/cmd/loggers"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"github.com/CyrilSbrodov/metricService.git/pkg/client/postgresql"
)

// PGSStore объявление структуры PostgreSQL.
type PGSStore struct {
	client postgresql.Client
	Hash   string
	logger loggers.Logger
}

// Создание таблиц в базе.
func createTable(ctx context.Context, client postgresql.Client, logger *loggers.Logger) error {
	tx, err := client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.LogErr(err, "failed to begin transaction")
		return err
	}
	defer tx.Rollback(ctx)

	q := `CREATE TABLE if not exists metrics (
    id TEXT NOT NULL,
    mType TEXT NOT NULL,
    delta BIGINT,
    value DOUBLE PRECISION,
    hash TEXT
	); 
	CREATE UNIQUE INDEX if not exists metrics_id_uindex on metrics (id);`

	_, err = tx.Exec(ctx, q)
	if err != nil {
		logger.LogErr(err, "failed to create table")
		return err
	}

	return tx.Commit(ctx)
}

// NewPGSStore создание новой базы.
func NewPGSStore(client postgresql.Client, cfg *config.ServerConfig, logger *loggers.Logger) (*PGSStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := createTable(ctx, client, logger); err != nil {
		logger.LogErr(err, "failed to create table")
		return nil, err
	}

	return &PGSStore{
		client: client,
		Hash:   cfg.Hash,
	}, nil
}

// GetMetric выгрузка метрики.
func (p *PGSStore) GetMetric(metric storage.Metrics) (storage.Metrics, error) {
	var m storage.Metrics
	if metric.MType == "counter" || metric.MType == "gauge" {
		q := `SELECT id, mType, delta, value, hash FROM metrics WHERE id = $1 AND mType = $2`
		if err := p.client.QueryRow(context.Background(), q, metric.ID, metric.MType).Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				p.logger.LogErr(err, "Failure to select object from table")
				return m, err
			}
			p.logger.LogErr(err, "wrong metric")
			return m, fmt.Errorf("missing metric %s", metric.ID)
		}
		m.Hash, _ = hashing(p.Hash, &m, &p.logger)
		return m, nil
	} else {
		p.logger.LogErr(fmt.Errorf("wrong type"), "wrong type")
		return m, fmt.Errorf("wrong type")
	}
}

// GetAll выгрузка всех метрик.
func (p *PGSStore) GetAll() (string, error) {
	var metrics []storage.Metrics
	q := `SELECT id, mType, delta, value, hash FROM metrics`
	rows, err := p.client.Query(context.Background(), q)
	if err != nil {
		p.logger.LogErr(err, "Failure to select object from table")
		return "", err
	}
	for rows.Next() {
		var m storage.Metrics
		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash)
		if err != nil {
			p.logger.LogErr(err, "Failure to convert object from table")
			return "", err
		}
		metrics = append(metrics, m)
	}

	result := ""
	for _, f := range metrics {
		if f.MType == "gauge" {
			if f.Value != nil {
				result += fmt.Sprintf("%s : %f\n", f.ID, *f.Value)
			}
			continue
		} else if f.MType == "counter" {
			result += fmt.Sprintf("%s : %d\n", f.ID, *f.Delta)
		}
	}
	return result, nil
}

// CollectMetric Сохранение метрики в БД.
func (p *PGSStore) CollectMetric(m storage.Metrics) error {

	if m.Hash != "" {
		_, ok := hashing(p.Hash, &m, &p.logger)
		if !ok {
			err := fmt.Errorf("hash is wrong")
			p.logger.LogErr(err, "hash is wrong")
			return err
		}
	}
	if m.MType == "counter" || m.MType == "gauge" {
		q := `INSERT INTO metrics (id, mType, delta, value, hash)
    						VALUES ($1, $2, $3, $4, $5)
							ON CONFLICT (id) DO UPDATE SET
    							delta = metrics.delta + EXCLUDED.delta,
    							value = $4,
    							hash = EXCLUDED.hash`
		if _, err := p.client.Exec(context.Background(), q, m.ID, m.MType, m.Delta, m.Value, m.Hash); err != nil {
			p.logger.LogErr(err, "Failure to insert object into table")
			return err
		}
	} else {
		p.logger.LogErr(fmt.Errorf("wrong type"), "wrong type")
		return fmt.Errorf("wrong type")
	}
	return nil
}

// CollectOrChangeGauge Сохранение или изменение метрики типа Gauge.
func (p *PGSStore) CollectOrChangeGauge(id string, value float64) error {
	mType := "gauge"
	hash := ""
	q := `INSERT INTO metrics (id, mType, value, hash)
    						VALUES ($1, $2, $3, $4)
							ON CONFLICT (id) DO UPDATE SET
    							value = EXCLUDED.value,
							    mType = EXCLUDED.mType`
	if _, err := p.client.Exec(context.Background(), q, id, mType, value, hash); err != nil {
		p.logger.LogErr(err, "Failure to insert object into table")
		return err
	}
	return nil
}

// CollectOrIncreaseCounter Сохранение или изменение метрики типа Counter.
func (p *PGSStore) CollectOrIncreaseCounter(id string, delta int64) error {
	mType := "counter"
	hash := ""
	q := `INSERT INTO metrics (id, mType, delta, hash)
    						VALUES ($1, $2, $3, $4)
							ON CONFLICT (id) DO UPDATE SET
    							delta = metrics.delta + EXCLUDED.delta,
 								mType = EXCLUDED.mType`
	if _, err := p.client.Exec(context.Background(), q, id, mType, delta, hash); err != nil {
		p.logger.LogErr(err, "Failure to insert object into table")
		return err
	}
	return nil
}

// GetGauge Выгрузка метрики типа Gauge.
func (p *PGSStore) GetGauge(id string) (float64, error) {
	var value float64
	mType := "gauge"
	q := `SELECT value FROM metrics WHERE id = $1 AND mType = $2`
	if err := p.client.QueryRow(context.Background(), q, id, mType).Scan(&value); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.logger.LogErr(err, "Failure to select object from table")
			return 0, err
		}
		p.logger.LogErr(err, "Wrong metric")
		return 0, fmt.Errorf("missing metric %s", id)
	}
	return value, nil
}

// GetCounter Выгрузка метрики типа Counter.
func (p *PGSStore) GetCounter(id string) (int64, error) {
	var delta int64
	mType := "counter"
	q := `SELECT delta FROM metrics WHERE id = $1 AND mType = $2`
	if err := p.client.QueryRow(context.Background(), q, id, mType).Scan(&delta); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.logger.LogErr(err, "Failure to select object from table")
			return 0, err
		}
		p.logger.LogErr(err, "Wrong metric")
		return 0, fmt.Errorf("missing metric %s", id)
	}
	return delta, nil
}

func (p *PGSStore) PingClient() error {
	return p.client.Ping(context.Background())
}

// CollectMetrics Сохранение метрики батчами.
func (p *PGSStore) CollectMetrics(metrics []storage.Metrics) error {

	tx, err := p.client.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		p.logger.LogErr(err, "failed to begin transaction")
		return err
	}
	defer tx.Rollback(context.Background())
	q := `INSERT INTO metrics (id, mType, delta, value, hash)
    						VALUES ($1, $2, $3, $4, $5)
							ON CONFLICT (id) DO UPDATE SET
    							delta = metrics.delta + EXCLUDED.delta,
    							value = $4,
    							hash = EXCLUDED.hash`

	for _, m := range metrics {
		if _, err = tx.Exec(context.Background(), q, m.ID, m.MType, m.Delta, m.Value, m.Hash); err != nil {
			p.logger.LogErr(err, "failed transaction")
			return err
		}
	}

	return tx.Commit(context.Background())
}
