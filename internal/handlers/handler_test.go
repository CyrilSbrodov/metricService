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

	"github.com/CyrilSbrodov/metricService.git/internal/handlers"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"github.com/CyrilSbrodov/metricService.git/internal/storage/repositories"
)

func TestHandler_GaugeHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	repo := repositories.NewRepository()

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
	repo := repositories.NewRepository()

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
	repo := repositories.NewRepository()

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

func TestHandler_CollectHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	var value float64 = 123123

	repo := repositories.NewRepository()

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
