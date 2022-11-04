package repositories

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"github.com/CyrilSbrodov/metricService.git/pkg/client/postgresql"
)

type PGSStore struct {
	client postgresql.Client
	Hash   string
}

func createTable(ctx context.Context, client postgresql.Client) error {
	tx, err := client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
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
		return err
	}
	return tx.Commit(ctx)
}

func NewPGSStore(client postgresql.Client, cfg *config.ServerConfig) (*PGSStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := createTable(ctx, client); err != nil {
		return nil, err
	}

	return &PGSStore{
		client: client,
		Hash:   cfg.Hash,
	}, nil
}

func (p *PGSStore) GetMetric(metric storage.Metrics) (storage.Metrics, error) {
	var m storage.Metrics
	if metric.MType == "counter" || metric.MType == "gauge" {
		q := `SELECT id, mType, delta, value, hash FROM metrics WHERE id = $1 AND mType = $2`
		if err := p.client.QueryRow(context.Background(), q, metric.ID, metric.MType).Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash); err != nil {
			if err != pgx.ErrNoRows {
				log.Fatal("Failure to select object from table. Due error: ", err)
				return m, err
			}
			return m, fmt.Errorf("missing metric %s", metric.ID)
		}
		m.Hash, _ = hashing(p.Hash, &m)
		return m, nil
	} else {
		return m, fmt.Errorf("wrong type")
	}
}

func (p *PGSStore) GetAll() (string, error) {
	var metrics []storage.Metrics
	q := `SELECT id, mType, delta, value, hash FROM metrics`
	rows, err := p.client.Query(context.Background(), q)
	if err != nil {
		log.Fatal("Failure to select object from table. Due error: ", err)
		return "", err
	}
	for rows.Next() {
		var m storage.Metrics
		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash)
		if err != nil {
			log.Fatal("Failure to convert object from table. Due error: ", err)
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

func (p *PGSStore) CollectMetrics(m storage.Metrics) error {

	if m.Hash != "" {
		_, ok := hashing(p.Hash, &m)
		if !ok {
			err := fmt.Errorf("hash is wrong")
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
			log.Fatal("Failure to insert object into table. Due error: ", err)
			return err
		}
	} else {
		return fmt.Errorf("wrong type")
	}
	return nil
}

func (p *PGSStore) CollectOrChangeGauge(id string, value float64) error {
	mType := "gauge"
	hash := ""
	q := `INSERT INTO metrics (id, mType, value, hash)
    						VALUES ($1, $2, $3, $4)
							ON CONFLICT (id) DO UPDATE SET
    							value = EXCLUDED.value,
							    mType = EXCLUDED.mType`
	if _, err := p.client.Exec(context.Background(), q, id, mType, value, hash); err != nil {
		log.Fatal("Failure to insert object into table. Due error: ", err)
		return err
	}
	return nil
}

func (p *PGSStore) CollectOrIncreaseCounter(id string, delta int64) error {
	mType := "counter"
	hash := ""
	q := `INSERT INTO metrics (id, mType, delta, hash)
    						VALUES ($1, $2, $3, $4)
							ON CONFLICT (id) DO UPDATE SET
    							delta = metrics.delta + EXCLUDED.delta,
 								mType = EXCLUDED.mType`
	if _, err := p.client.Exec(context.Background(), q, id, mType, delta, hash); err != nil {
		log.Fatal("Failure to insert object into table. Due error: ", err)
		return err
	}
	return nil
}

func (p *PGSStore) GetGauge(id string) (float64, error) {
	var value float64
	mType := "gauge"
	q := `SELECT value FROM metrics WHERE id = $1 AND mType = $2`
	if err := p.client.QueryRow(context.Background(), q, id, mType).Scan(&value); err != nil {
		if err != pgx.ErrNoRows {
			log.Fatal("Failure to select object from table. Due error: ", err)
			return 0, err
		}
		return 0, fmt.Errorf("missing metric %s", id)
	}
	return value, nil
}

func (p *PGSStore) GetCounter(id string) (int64, error) {
	var delta int64
	mType := "counter"
	q := `SELECT delta FROM metrics WHERE id = $1 AND mType = $2`
	if err := p.client.QueryRow(context.Background(), q, id, mType).Scan(&delta); err != nil {
		if err != pgx.ErrNoRows {
			log.Fatal("Failure to select object from table. Due error: ", err)
			return 0, err
		}
		return 0, fmt.Errorf("missing metric %s", id)
	}
	return delta, nil
}

func (p *PGSStore) PingClient() error {
	return p.client.Ping(context.Background())
}

//функция загрузки данных на диск
func (p *PGSStore) Upload() error {
	return nil
}

// TODO рещить проблему, чтобы файл был только для repository.
func (p *PGSStore) UploadWithTicker(ticker *time.Ticker, done chan os.Signal) {
	for {
		select {
		case <-ticker.C:
			err := p.Upload()
			if err != nil {
				fmt.Println(err)
				return
			}
		case <-done:
			ticker.Stop()
			return
		}
	}
}
