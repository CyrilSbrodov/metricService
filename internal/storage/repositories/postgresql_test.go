package repositories

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

var (
	CFG config.ServerConfig
)

func TestMain(m *testing.M) {
	cfg := config.ServerConfigInit()
	CFG = *cfg
	CFG.DatabaseDSN = "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"
	os.Exit(m.Run())
}

func TestPGSStore_CollectMetric(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("metrics")
	var value float64 = 100
	err := s.CollectMetric(storage.Metrics{
		ID:    "test",
		MType: "gauge",
		Delta: nil,
		Value: &value,
		Hash:  "",
	})
	assert.NoError(t, err)

}

func TestPGSStore_CollectMetrics(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("metrics")
	var value float64 = 100
	err := s.CollectMetrics([]storage.Metrics{
		{
			ID:    "test",
			MType: "gauge",
			Delta: nil,
			Value: &value,
			Hash:  "",
		},
	})
	assert.NoError(t, err)
}

func TestPGSStore_CollectOrChangeGauge(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("metrics")
	err := s.CollectOrChangeGauge("test", 100)
	assert.NoError(t, err)
}

func TestPGSStore_CollectOrIncreaseCounter(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("metrics")
	err := s.CollectOrIncreaseCounter("test", 100)
	assert.NoError(t, err)
}

func TestPGSStore_GetAll(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("metrics")
	err := s.CollectOrChangeGauge("testGauge", 100)
	assert.NoError(t, err)
	err = s.CollectOrIncreaseCounter("testCounter", 100)
	assert.NoError(t, err)
	m, err := s.GetAll()
	assert.NotNil(t, m)
	assert.NoError(t, err)
}

func TestPGSStore_GetCounter(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("metrics")
	err := s.CollectOrIncreaseCounter("test", 100)
	assert.NoError(t, err)
	m, err := s.GetCounter("test")
	assert.NotNil(t, m)
	assert.NoError(t, err)
}

func TestPGSStore_GetGauge(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("metrics")
	err := s.CollectOrChangeGauge("test", 100)
	assert.NoError(t, err)
	m, err := s.GetGauge("test")
	assert.NotNil(t, m)
	assert.NoError(t, err)
}

func TestPGSStore_GetMetric(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("metrics")
	var value float64 = 100

	err := s.CollectMetrics([]storage.Metrics{
		{
			ID:    "test",
			MType: "gauge",
			Delta: nil,
			Value: &value,
			Hash:  "",
		},
	})
	assert.NoError(t, err)

	m, err := s.GetMetric(storage.Metrics{
		ID:    "test",
		MType: "gauge",
		Delta: nil,
		Value: &value,
		Hash:  "",
	})
	assert.NoError(t, err)
	assert.NotNil(t, m)
}

func TestPGSStore_PingClient(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("metrics")
	err := s.PingClient()
	assert.NoError(t, err)
}
