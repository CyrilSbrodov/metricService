package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/cmd/loggers"
	"github.com/CyrilSbrodov/metricService.git/internal/handlers"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"github.com/CyrilSbrodov/metricService.git/internal/storage/repositories"
)

var (
	CFG config.ServerConfig
)

func TestMain(m *testing.M) {
	cfg := config.ServerConfigInit()
	CFG = cfg
	os.Exit(m.Run())
}

func TestHandler_CollectHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	var value float64 = 123123
	logger := loggers.NewLogger()
	repo, _ := repositories.NewRepository(&CFG, logger)

	type fields struct {
		Storage storage.Storage
	}
	tests := []struct {
		name    string
		fields  fields
		want    want
		request string
		req     storage.Metrics
	}{
		{
			name: "Test ok",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update",
			want: want{
				200,
			},
			req: storage.Metrics{
				ID:    "Alloc",
				MType: "gauge",
				Value: &value,
			},
		},
		{
			name: "Test ok",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/",
			want: want{
				200,
			},
			req: storage.Metrics{
				ID:    "Alloc",
				MType: "gauge",
				Value: &value,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricsJSON, err := json.Marshal(tt.req)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			request := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewBuffer(metricsJSON))
			w := httptest.NewRecorder()
			h := handlers.Handler{
				Storage: tt.fields.Storage,
			}
			h.CollectHandler().ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_GetAllHandler(t *testing.T) {
	type want struct {
		statusCode int
	}

	logger := loggers.NewLogger()
	repo, _ := repositories.NewRepository(&CFG, logger)

	type fields struct {
		Storage storage.Storage
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "status 200",
			fields: fields{
				Storage: repo,
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			h := handlers.Handler{
				Storage: tt.fields.Storage,
			}
			h.GetAllHandler().ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_GaugeHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	logger := loggers.NewLogger()
	repo, _ := repositories.NewRepository(&CFG, logger)

	type fields struct {
		Storage storage.Storage
	}
	tests := []struct {
		name    string
		request string
		fields  fields
		want    want
	}{
		{
			name: "Test ok",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/gauge/test/100",
			want: want{
				200,
			},
		},
		{
			name: "Test 501",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/gaug/test/100",
			want: want{
				501,
			},
		},
		{
			name: "Test 404",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/updat/gaug/test/100",
			want: want{
				404,
			},
		},
		{
			name: "Test 400",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/gauge/test/none",
			want: want{
				400,
			},
		},
		{
			name: "Test 404",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/gauge/test",
			want: want{
				404,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := handlers.Handler{
				Storage: tt.fields.Storage,
			}
			h.GaugeHandler().ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_CounterHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	logger := loggers.NewLogger()
	repo, _ := repositories.NewRepository(&CFG, logger)

	type fields struct {
		Storage storage.Storage
	}
	tests := []struct {
		name    string
		request string
		fields  fields
		want    want
	}{
		{
			name: "Test ok",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/counter/test/100",
			want: want{
				200,
			},
		},
		{
			name: "Test 501",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/test/test/100",
			want: want{
				501,
			},
		},
		{
			name: "Test 400",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/counter/test/none",
			want: want{
				400,
			},
		},
		{
			name: "Test 404",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/counter/test",
			want: want{
				404,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := handlers.Handler{
				Storage: tt.fields.Storage,
			}
			h.CounterHandler().ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_OtherHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	logger := loggers.NewLogger()
	repo, _ := repositories.NewRepository(&CFG, logger)

	type fields struct {
		Storage storage.Storage
	}
	tests := []struct {
		name    string
		request string
		fields  fields
		want    want
	}{
		{
			name: "Test wrong path/method, code 404",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/test/counter/test/100",
			want: want{
				404,
			},
		},
		{
			name: "Test wrong types, code 501",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/test/test/100",
			want: want{
				501,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := handlers.Handler{
				Storage: tt.fields.Storage,
			}
			h.OtherHandler().ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_GetHandlerJSON(t *testing.T) {
	var delta int64 = 100
	logger := loggers.NewLogger()
	repo, _ := repositories.NewRepository(&CFG, logger)

	repo.Metrics = make(map[string]storage.Metrics)
	var m = storage.Metrics{
		ID:    "test",
		MType: "counter",
		Delta: &delta,
	}
	repo.Metrics[m.ID] = m

	type fields struct {
		Storage storage.Storage
		logger  loggers.Logger
	}
	tests := []struct {
		name         string
		body         storage.Metrics
		expectedCode int
		fields       fields
	}{
		{
			name: "Test ok",
			body: m,
			fields: fields{
				Storage: repo,
				logger:  *logger,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 404",
			body: storage.Metrics{
				ID:    "testWrong",
				MType: "counter",
				Delta: &delta,
			},
			fields: fields{
				Storage: repo,
				logger:  *logger,
			},
			expectedCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricsJSON, err := json.Marshal(tt.body)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			request := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(metricsJSON))
			w := httptest.NewRecorder()
			h := handlers.Handler{
				Storage: tt.fields.Storage,
			}
			h.GetHandlerJSON().ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.expectedCode, result.StatusCode)
		})
	}
}

func TestHandler_GetHandler(t *testing.T) {
	var delta int64 = 100
	var val float64 = 100
	logger := loggers.NewLogger()
	repo, _ := repositories.NewRepository(&CFG, logger)

	repo.Metrics = make(map[string]storage.Metrics)
	var m = storage.Metrics{
		ID:    "test",
		MType: "gauge",
		Value: &val,
	}
	repo.Counter[m.ID] = delta
	repo.Metrics[m.ID] = m
	type fields struct {
		Storage storage.Storage
		logger  loggers.Logger
	}
	tests := []struct {
		name         string
		request      string
		expectedCode int
		fields       fields
	}{
		{
			name:    "Test ok",
			request: "http://localhost:8080/value/counter/test/100",
			fields: fields{
				Storage: repo,
				logger:  *logger,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:    "Test 404",
			request: "http://localhost:8080/value/gauge/testWrong/100",
			fields: fields{
				Storage: repo,
				logger:  *logger,
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "Test 501",
			request: "http://localhost:8080/value/test/100",
			fields: fields{
				Storage: repo,
				logger:  *logger,
			},
			expectedCode: http.StatusNotImplemented,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := handlers.Handler{
				Storage: tt.fields.Storage,
			}
			h.GetHandler().ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.expectedCode, result.StatusCode)
		})
	}
}

func TestHandler_CollectBatchHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	var value float64 = 123123
	logger := loggers.NewLogger()
	repo, _ := repositories.NewRepository(&CFG, logger)

	type fields struct {
		Storage storage.Storage
	}

	tests := []struct {
		name    string
		fields  fields
		want    want
		request string
		req     []storage.Metrics
	}{
		{
			name: "Test ok",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/updates",
			want: want{
				200,
			},
			req: []storage.Metrics{
				{
					ID:    "Alloc",
					MType: "gauge",
					Value: &value,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricsJSON, err := json.Marshal(tt.req)
			assert.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewBuffer(metricsJSON))
			w := httptest.NewRecorder()
			h := handlers.Handler{
				Storage: tt.fields.Storage,
			}
			h.CollectBatchHandler().ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func newRepo() *repositories.Repository {
	logger := loggers.NewLogger()
	repo, _ := repositories.NewRepository(&CFG, logger)
	return repo
}

func newMetric() storage.Metrics {
	var value float64 = 100
	return storage.Metrics{
		ID:    "test",
		MType: "gauge",
		Delta: nil,
		Value: &value,
		Hash:  "",
	}
}
func ExampleHandler_CollectHandler() {
	repo := newRepo()
	m := newMetric()

	metricsJSON, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req := "http://localhost:8080/update/"
	request := httptest.NewRequest(http.MethodPost, req, bytes.NewBuffer(metricsJSON))
	w := httptest.NewRecorder()
	h := handlers.Handler{
		Storage: repo,
	}
	h.CollectHandler().ServeHTTP(w, request)
	result := w.Result()
	defer result.Body.Close()
	fmt.Println(result.StatusCode)

	//Output:
	//200
}

func ExampleHandler_GetAllHandler() {
	repo := newRepo()
	var value float64 = 150
	repo.Metrics["test"] = storage.Metrics{
		ID:    "test",
		MType: "gauge",
		Delta: nil,
		Value: &value,
		Hash:  "",
	}
	repo.Metrics["Alloc"] = newMetric()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h := handlers.Handler{
		Storage: repo,
	}
	h.GetAllHandler().ServeHTTP(w, request)
	metrics, err := h.GetAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	result := w.Result()
	defer result.Body.Close()
	fmt.Println(result.StatusCode)
	fmt.Println(metrics)

	//Output:
	//200
	//Alloc : 100.000000<br>
	//test : 150.000000<br>

}
