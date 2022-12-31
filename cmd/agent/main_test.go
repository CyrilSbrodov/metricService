package main

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

var (
	CFG config.AgentConfig
)

func TestMain(m *testing.M) {
	cfg := config.AgentConfigInit()
	CFG = cfg
	os.Exit(m.Run())
}

func Test_update(t *testing.T) {
	store := storage.NewAgentMetrics()
	wg := &sync.WaitGroup{}
	type args struct {
		store *storage.AgentMetrics
		count int64
		cfg   *config.AgentConfig
		wg    *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test OK",
			args: args{
				store: store,
				count: 5,
				cfg:   &CFG,
				wg:    wg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.wg.Add(1)
			update(tt.args.store, tt.args.count, tt.args.cfg, tt.args.wg)
			assert.NotNil(t, store)
		})
	}
}

func Test_updateOtherMetrics(t *testing.T) {
	store := storage.NewAgentMetrics()
	wg := &sync.WaitGroup{}
	type args struct {
		store *storage.AgentMetrics
		cfg   *config.AgentConfig
		wg    *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test OK",
			args: args{
				store: store,
				cfg:   &CFG,
				wg:    wg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.wg.Add(1)
			updateOtherMetrics(tt.args.store, tt.args.wg, tt.args.cfg)
			assert.NotNil(t, store)
		})
	}
}

func BenchmarkUpdate(b *testing.B) {
	store := storage.NewAgentMetrics()
	wg := &sync.WaitGroup{}
	var count int64 = 0
	b.Run("test", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			update(store, count, &CFG, wg)
		}
	})
}
