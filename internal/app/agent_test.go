package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/cmd/loggers"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

var (
	CFG config.AgentConfig
)

func TestMain(m *testing.M) {
	cfg := config.AgentConfigInit()
	CFG = *cfg
	os.Exit(m.Run())
}

//type AgentTestSuite struct {
//	suite.Suite
//	AgentApp AgentApp
//}
//
//func (suite *AgentTestSuite) SetupTest() {
//	suite.AgentApp.client = http.DefaultClient
//	suite.AgentApp.cfg = config.AgentConfigInit()
//	suite.AgentApp.logger = loggers.NewLogger()
//
//}
//
//func (suite *AgentTestSuite) TestAgent() {
//	suite.Run("receive data from agent", func() {
//	})
//}

//func TestAgentApp_Run(t *testing.T) {
//	logger := loggers.NewLogger()
//	cfg := config.AgentConfigInit()
//	type fields struct {
//		client             *http.Client
//		cfg                config.AgentConfig
//		logger             *loggers.Logger
//		expectedRespStatus int
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		{
//			name: "Run",
//			fields: fields{
//				client:             http.DefaultClient,
//				cfg:                cfg,
//				logger:             logger,
//				expectedRespStatus: http.StatusOK,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			http.DefaultServeMux.ServeHTTP(resp, req)
//			if p, err := ioutil.ReadAll(resp.Body); err != nil {
//				t.Fail()
//			} else {
//				if strings.Contains(string(p), "Error") {
//					t.Errorf("header response shouldn't return error: %s", p)
//				} else if !strings.Contains(string(p), `expected result`) {
//					t.Errorf("header response doen't match:\n%s", p)
//				}
//			}
//
//			a := &AgentApp{
//				client: tt.fields.client,
//				cfg:    tt.fields.cfg,
//				logger: tt.fields.logger,
//			}
//			a.Run()
//		})
//	}
//}

func TestAgentApp_compress(t *testing.T) {
	var b = []byte{1}
	type fields struct {
		client *http.Client
		cfg    config.AgentConfig
		logger *loggers.Logger
	}
	type args struct {
		store []byte
		count int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			fields: fields{
				client: nil,
				cfg:    config.AgentConfig{},
				logger: nil,
			},
			args: args{
				store: b,
			},
			want:    []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 98, 4, 4, 0, 0, 255, 255, 27, 223, 5, 165, 1, 0, 0, 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AgentApp{
				client: tt.fields.client,
				cfg:    tt.fields.cfg,
				logger: tt.fields.logger,
			}

			got, err := a.compress(tt.args.store)
			if (err != nil) != tt.wantErr {
				t.Errorf("compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgentApp_hashing(t *testing.T) {
	var value float64 = 1
	var m = storage.Metrics{
		ID:    "1",
		MType: "gauge",
		Delta: nil,
		Value: &value,
		Hash:  "",
	}
	type fields struct {
		client *http.Client
		cfg    config.AgentConfig
		logger *loggers.Logger
	}
	type args struct {
		m *storage.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{
				client: nil,
				cfg:    config.AgentConfig{},
				logger: nil,
			},
			args: args{
				m: &m,
			},

			want: "c1da6a11da46ad4466e93330a3e06437846b371de51cdda9baf84d862ccce1b0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AgentApp{
				client: tt.fields.client,
				cfg:    tt.fields.cfg,
				logger: tt.fields.logger,
			}
			if got := a.hashing(tt.args.m); got != tt.want {
				t.Errorf("hashing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgentApp_update(t *testing.T) {
	store := storage.NewAgentMetrics()
	wg := &sync.WaitGroup{}
	type fields struct {
		client *http.Client
		cfg    config.AgentConfig
		logger *loggers.Logger
	}
	type args struct {
		store *storage.AgentMetrics
		count int64
		wg    *sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			fields: fields{
				client: nil,
				cfg:    config.AgentConfig{},
				logger: nil,
			},
			args: args{
				store: store,
				count: 0,
				wg:    wg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AgentApp{
				client: tt.fields.client,
				cfg:    tt.fields.cfg,
				logger: tt.fields.logger,
			}
			tt.args.wg.Add(1)
			a.update(tt.args.store, tt.args.count, tt.args.wg)
			assert.NotNil(t, tt.args.store)
		})
	}
}

func TestAgentApp_updateOtherMetrics(t *testing.T) {
	store := storage.NewAgentMetrics()
	wg := &sync.WaitGroup{}
	type fields struct {
		client *http.Client
		cfg    config.AgentConfig
		logger *loggers.Logger
	}
	type args struct {
		store *storage.AgentMetrics
		count int64
		wg    *sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			fields: fields{
				client: nil,
				cfg:    config.AgentConfig{},
				logger: nil,
			},
			args: args{
				store: store,
				count: 0,
				wg:    wg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AgentApp{
				client: tt.fields.client,
				cfg:    tt.fields.cfg,
				logger: tt.fields.logger,
			}
			tt.args.wg.Add(1)
			a.updateOtherMetrics(tt.args.store, tt.args.wg)
			assert.NotNil(t, tt.args.store)
		})
	}
}

func TestAgentApp_uploadBatch(t *testing.T) {
	client := http.DefaultClient
	store := storage.NewAgentMetrics()
	logger := loggers.NewLogger()
	wg := &sync.WaitGroup{}
	var delta int64 = 1
	var m = storage.Metrics{
		ID:    "1",
		MType: "counter",
		Delta: &delta,
		Value: nil,
		Hash:  "",
	}
	store.Store[m.ID] = m
	type fields struct {
		client *http.Client
		cfg    config.AgentConfig
		logger *loggers.Logger
	}
	type args struct {
		store *storage.AgentMetrics
		wg    *sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			fields: fields{
				client: client,
				cfg:    CFG,
				logger: logger,
			},
			args: args{
				store: store,
				wg:    wg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"ok"}`)
			}))
			tt.fields.cfg.Addr = ts.URL
			defer ts.Close()
			a := &AgentApp{
				client: tt.fields.client,
				cfg:    tt.fields.cfg,
				logger: tt.fields.logger,
			}
			tt.args.wg.Add(1)
			a.uploadBatch(tt.args.store, tt.args.wg)
		})
	}
}

func TestAgentApp_upload(t *testing.T) {
	client := http.DefaultClient
	store := storage.NewAgentMetrics()
	logger := loggers.NewLogger()
	wg := &sync.WaitGroup{}
	var delta int64 = 1
	var m = storage.Metrics{
		ID:    "1",
		MType: "counter",
		Delta: &delta,
		Value: nil,
		Hash:  "",
	}
	store.Store[m.ID] = m
	type fields struct {
		client *http.Client
		cfg    config.AgentConfig
		logger *loggers.Logger
	}
	type args struct {
		store *storage.AgentMetrics
		wg    *sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			fields: fields{
				client: client,
				cfg:    CFG,
				logger: logger,
			},
			args: args{
				store: store,
				wg:    wg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"ok"}`)
			}))
			tt.fields.cfg.Addr = ts.URL
			defer ts.Close()
			a := &AgentApp{
				client: tt.fields.client,
				cfg:    tt.fields.cfg,
				logger: tt.fields.logger,
			}
			tt.args.wg.Add(1)
			a.upload(tt.args.store, tt.args.wg)
		})
	}
}
