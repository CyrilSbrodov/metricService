package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"github.com/CyrilSbrodov/metricService.git/pkg/client/postgresql"
)

type PGSStore struct {
	client  postgresql.Client
	Metrics map[string]storage.Metrics
	//Dsn           string
	//file          *os.File
	//Check         bool
	//StoreInterval time.Duration
	Hash string
}

func createTable(ctx context.Context, client postgresql.Client) error {
	tx, err := client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	q := `CREATE TABLE if not exists metrics (
    id TEXT NOT NULL,
    mType TEXT NOT NULL,
    delta INT,
    value DOUBLE PRECISION,
    hash TEXT
	); 
	CREATE UNIQUE INDEX if not exists metrics_id_uindex on metrics (id);`
	_, err = tx.Exec(ctx, q)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func NewPGSStore(client postgresql.Client, cfg *config.ServerConfig) (*PGSStore, error) {
	metrics := storage.MetricsStore

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := createTable(ctx, client); err != nil {
		return nil, err
	}
	return &PGSStore{
		Metrics: metrics,
		client:  client,
		Hash:    cfg.Hash,
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

	if m.Hash != "" {
		_, ok := hashing(p.Hash, &m)
		if !ok {
			err := fmt.Errorf("hash is wrong")
			return err
		}
	}
	var metric storage.Metrics
	if m.MType == "counter" || m.MType == "gauge" {
		sqlStatement := `INSERT INTO metrics (id, mType, delta, value, hash)
    						VALUES ($1, $2, $3, $4, $5)
							ON CONFLICT (id) DO UPDATE SET
    							delta = metrics.delta + EXCLUDED.delta,
    							value = $4,
    							hash = $5
    							RETURNING id`
		if err := p.client.QueryRow(context.Background(), sqlStatement, m.ID, m.MType, m.Delta, m.Value, m.Hash).Scan(&metric.ID); err != nil {
			log.Fatal("Ошибка добавления данных в БД. ", err)
		}
		fmt.Println(metric)
	} else {
		return fmt.Errorf("wrong type")
	}
	return nil
}

func (p *PGSStore) CollectOrChangeGauge(id string, value float64) error {
	mType := "gauge"
	sqlStatement := `INSERT INTO metrics (id, mType, value)
    						VALUES ($1, $2, $3)
							ON CONFLICT (id) DO UPDATE SET
    							value = EXCLUDED.value,
							    mType = EXCLUDED.mType`
	if err := p.client.QueryRow(context.Background(), sqlStatement, id, value, mType); err != nil {
		log.Fatal("Ошибка добавления данных в БД. ", err)
	}
	return nil
}

func (p *PGSStore) CollectOrIncreaseCounter(id string, delta int64) error {
	mType := "counter"
	sqlStatement := `INSERT INTO metrics (id, mType, delta)
    						VALUES ($1, $2, $3)
							ON CONFLICT (id) DO UPDATE SET
    							delta = metrics.delta + EXCLUDED.delta,
 								mType = EXCLUDED.mType`
	if err := p.client.QueryRow(context.Background(), sqlStatement, id, delta, mType); err != nil {
		log.Fatal("Ошибка добавления данных в БД. ", err)
	}
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
