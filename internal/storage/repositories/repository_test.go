package repositories

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

func TestRepository_CollectMetric(t *testing.T) {
	repo, _ := NewRepository(&CFG, nil)
	type fields struct {
		repo *Repository
	}
	var delta int64 = 100
	var value float64 = 100
	type args struct {
		m storage.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test counter OK",
			fields: fields{
				repo: repo,
			},
			args: args{m: storage.Metrics{
				ID:    "test",
				MType: "counter",
				Delta: &delta,
				Value: nil,
				Hash:  "",
			}},
		},
		{
			name: "test gauge OK",
			fields: fields{
				repo: repo,
			},
			args: args{m: storage.Metrics{
				ID:    "test",
				MType: "gauge",
				Delta: nil,
				Value: &value,
				Hash:  "",
			}},
		},
		{
			name: "test other MType OK",
			fields: fields{
				repo: repo,
			},
			args: args{m: storage.Metrics{
				ID:    "test",
				MType: "test",
				Delta: nil,
				Value: &value,
				Hash:  "",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.repo.CollectMetric(tt.args.m)
			assert.NoError(t, err)
		})
	}
}

func TestRepository_CollectMetrics(t *testing.T) {
	repo, _ := NewRepository(&CFG, nil)
	repo.Metrics = nil
	repo.Metrics = make(map[string]storage.Metrics)
	type fields struct {
		repo *Repository
	}
	var delta int64 = 100
	var value float64 = 100
	type args struct {
		m []storage.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test counter OK",
			fields: fields{
				repo: repo,
			},
			args: args{m: []storage.Metrics{
				{
					ID:    "test",
					MType: "counter",
					Delta: &delta,
					Value: nil,
					Hash:  "",
				},
			}},
		},
		{
			name: "test gauge OK",
			fields: fields{
				repo: repo,
			},
			args: args{m: []storage.Metrics{
				{
					ID:    "test",
					MType: "gauge",
					Delta: nil,
					Value: &value,
					Hash:  "",
				},
			}},
		},
		{
			name: "test other MType OK",
			fields: fields{
				repo: repo,
			},
			args: args{m: []storage.Metrics{
				{
					ID:    "test",
					MType: "test",
					Delta: nil,
					Value: &value,
					Hash:  "",
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.repo.CollectMetrics(tt.args.m)
			assert.NoError(t, err)
		})
	}
}

func TestRepository_CollectOrChangeGauge(t *testing.T) {
	repo, _ := NewRepository(&CFG, nil)
	type fields struct {
		repo *Repository
	}
	var value float64 = 100
	type args struct {
		m storage.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test OK",
			fields: fields{
				repo: repo,
			},
			args: args{m: storage.Metrics{
				ID:    "test",
				MType: "gauge",
				Delta: nil,
				Value: &value,
				Hash:  "",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.repo.CollectOrChangeGauge(tt.args.m.ID, *tt.args.m.Value)
			assert.NoError(t, err)
		})
	}
}

func TestRepository_CollectOrIncreaseCounter(t *testing.T) {
	repo, _ := NewRepository(&CFG, nil)
	type fields struct {
		repo *Repository
	}
	var delta int64 = 100
	type args struct {
		m storage.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test OK",
			fields: fields{
				repo: repo,
			},
			args: args{m: storage.Metrics{
				ID:    "test",
				MType: "counter",
				Delta: &delta,
				Value: nil,
				Hash:  "",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.repo.CollectOrIncreaseCounter(tt.args.m.ID, *tt.args.m.Delta)
			assert.NoError(t, err)
		})
	}
}

func TestRepository_GetAll(t *testing.T) {
	repo, _ := NewRepository(&CFG, nil)
	type fields struct {
		repo *Repository
	}
	var delta int64 = 100
	repo.Metrics["test"] = storage.Metrics{
		ID:    "test",
		MType: "counter",
		Delta: &delta,
		Value: nil,
		Hash:  "",
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test OK",
			fields: fields{
				repo: repo,
			},
			want: "test : 100<br>\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := tt.fields.repo.GetAll()
			assert.Equal(t, tt.want, m)
			assert.NoError(t, err)
		})
	}
}

func TestRepository_GetCounter(t *testing.T) {
	repo, _ := NewRepository(&CFG, nil)
	type fields struct {
		repo *Repository
	}
	var delta int64 = 100
	type args struct {
		m string
	}
	type want struct {
		delta int64
		value float64
		err   error
	}
	repo.Counter["test"] = delta
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "test OK",
			fields: fields{
				repo: repo,
			},
			args: args{
				m: "test",
			},
			want: want{
				delta: delta,
				err:   nil,
			},
		},
		{
			name: "test wrong id",
			fields: fields{
				repo: repo,
			},
			args: args{
				m: "test1",
			},
			want: want{
				delta: 0,
				err:   fmt.Errorf("missing metric test1"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := tt.fields.repo.GetCounter(tt.args.m)
			assert.Equal(t, tt.want.delta, m)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestRepository_GetGauge(t *testing.T) {
	repo, _ := NewRepository(&CFG, nil)
	type fields struct {
		repo *Repository
	}
	type want struct {
		delta int64
		value float64
		err   error
	}
	var value float64 = 100
	type args struct {
		m string
	}
	repo.Metrics["test"] = storage.Metrics{
		ID:    "test",
		MType: "gauge",
		Delta: nil,
		Value: &value,
		Hash:  "",
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "test OK",
			fields: fields{
				repo: repo,
			},
			args: args{
				m: "test",
			},
			want: want{
				value: value,
				err:   nil,
			},
		},
		{
			name: "test wrong id",
			fields: fields{
				repo: repo,
			},
			args: args{
				m: "test1",
			},
			want: want{
				value: 0,
				err:   fmt.Errorf("missing metric test1"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := tt.fields.repo.GetGauge(tt.args.m)
			assert.Equal(t, tt.want.value, m)
			assert.Equal(t, tt.want.err, err)

		})
	}
}

func TestRepository_GetMetric(t *testing.T) {
	repo, _ := NewRepository(&CFG, nil)
	type fields struct {
		repo *Repository
	}
	type want struct {
		m   storage.Metrics
		err error
	}
	var value float64 = 100
	var metric = storage.Metrics{
		ID:    "test",
		MType: "gauge",
		Delta: nil,
		Value: &value,
		Hash:  "",
	}
	type args struct {
		m storage.Metrics
	}
	hash, _ := hashing("", &metric, nil)
	repo.Metrics["test"] = metric
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "test OK",
			fields: fields{
				repo: repo,
			},
			args: args{
				m: storage.Metrics{
					ID:    "test",
					MType: "gauge",
					Delta: nil,
					Value: nil,
					Hash:  "",
				},
			},
			want: want{
				m: storage.Metrics{
					ID:    "test",
					MType: "gauge",
					Delta: nil,
					Value: &value,
					Hash:  hash,
				},
				err: nil,
			},
		},
		{
			name: "test wrong id",
			fields: fields{
				repo: repo,
			},
			args: args{
				m: storage.Metrics{
					ID:    "test1",
					MType: "gauge",
					Delta: nil,
					Value: nil,
					Hash:  "",
				},
			},
			want: want{
				m: storage.Metrics{
					ID:    "test1",
					MType: "gauge",
					Delta: nil,
					Value: nil,
					Hash:  "",
				},
				err: fmt.Errorf("id not found test1"),
			},
		},
		{
			name: "test wrong type",
			fields: fields{
				repo: repo,
			},
			args: args{
				m: storage.Metrics{
					ID:    "test1",
					MType: "test",
					Delta: nil,
					Value: nil,
					Hash:  "",
				},
			},
			want: want{
				m: storage.Metrics{
					ID:    "test1",
					MType: "test",
					Delta: nil,
					Value: nil,
					Hash:  "",
				},
				err: fmt.Errorf("type test is wrong"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := tt.fields.repo.GetMetric(tt.args.m)
			assert.Equal(t, tt.want.m, m)
			assert.Equal(t, tt.want.err, err)

		})
	}
}
