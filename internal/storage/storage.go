package storage

type Storage interface {
	GetMetric(metric Metrics) (Metrics, error)
	GetAll() string
	CollectMetrics(m Metrics) error
	CollectOrChangeGauge(name string, value float64) error
	CollectOrIncreaseCounter(name string, value int64) error
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
}

//type PostrgeStorage interface {
//	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
//	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
//	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
//	Begin(ctx context.Context) (pgx.Tx, error)
//	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
//	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
//}
