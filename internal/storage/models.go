package storage

import "sync"

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

// NewAgentMetrics инициализация нового временного хранилища для агента.
func NewAgentMetrics() *AgentMetrics {
	return &AgentMetrics{
		Store: make(map[string]Metrics),
	}
}
