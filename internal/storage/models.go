package storage

import (
	"sync"

	pb "github.com/CyrilSbrodov/metricService.git/internal/app/proto"
)

var GaugeData = map[string]float64{}
var CounterData = map[string]int64{}

// Metrics структура метрики.
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики.
	MType string   `json:"type"`            // параметр, принимающий значение Gauge или Counter.
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи Counter.
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи Gauge.
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции.
}

// MetricsStore инициализация временного хранилища.
var MetricsStore = map[string]Metrics{}

// AgentMetrics структура метрик для агента.
type AgentMetrics struct {
	Store map[string]Metrics
	Sync  sync.Mutex
}

type AgentMetricsProto struct {
	Store map[string]pb.Metrics
	Sync  sync.Mutex
}

// NewAgentMetrics инициализация нового временного хранилища для агента.
func NewAgentMetrics() *AgentMetrics {
	return &AgentMetrics{
		Store: make(map[string]Metrics),
	}
}

func NewAgentMetricsProto() *AgentMetricsProto {
	return &AgentMetricsProto{
		Store: make(map[string]pb.Metrics),
	}
}

func (m *Metrics) ToProto() *pb.Metrics {
	var value float64
	var delta int64
	if m.Value == nil {
		return &pb.Metrics{
			ID:    m.ID,
			MType: m.MType,
			Delta: *m.Delta,
			Value: value,
			Hash:  m.Hash,
		}
	} else {
		return &pb.Metrics{
			ID:    m.ID,
			MType: m.MType,
			Delta: delta,
			Value: *m.Value,
			Hash:  m.Hash,
		}
	}
}

func (m *Metrics) ToMetric(metric pb.Metrics) *Metrics {
	var value float64 = 0
	var delta int64 = 0
	if metric.Value == 0 {
		return &Metrics{
			ID:    metric.ID,
			MType: metric.MType,
			Delta: &metric.Delta,
			Value: &value,
			Hash:  metric.Hash,
		}
	} else {
		return &Metrics{
			ID:    metric.ID,
			MType: metric.MType,
			Delta: &delta,
			Value: &metric.Value,
			Hash:  metric.Hash,
		}
	}
}
