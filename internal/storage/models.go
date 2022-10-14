package storage

type User struct {
	ID        string
	FirstName string
	LastName  string
}

type Gauge struct {
	Alloc struct {
		Name  string
		Value float64
	}
	BuckHashSys struct {
		Name  string
		Value float64
	}
	Frees struct {
		Name  string
		Value float64
	}
	GCCPUFraction struct {
		Name  string
		Value float64
	}
	GCSys struct {
		Name  string
		Value float64
	}
	HeapAlloc struct {
		Name  string
		Value float64
	}
	HeapIdle struct {
		Name  string
		Value float64
	}
	HeapInuse struct {
		Name  string
		Value float64
	}
	HeapObjects struct {
		Name  string
		Value float64
	}
	HeapReleased struct {
		Name  string
		Value float64
	}
	HeapSys struct {
		Name  string
		Value float64
	}
	LastGC struct {
		Name  string
		Value float64
	}
	Lookups struct {
		Name  string
		Value float64
	}
	MCacheInuse struct {
		Name  string
		Value float64
	}
	MCacheSys struct {
		Name  string
		Value float64
	}
	MSpanInuse struct {
		Name  string
		Value float64
	}
	MSpanSys struct {
		Name  string
		Value float64
	}
	Mallocs struct {
		Name  string
		Value float64
	}
	NextGC struct {
		Name  string
		Value float64
	}
	NumForcedGC struct {
		Name  string
		Value float64
	}
	NumGC struct {
		Name  string
		Value float64
	}
	OtherSys struct {
		Name  string
		Value float64
	}
	PauseTotalNs struct {
		Name  string
		Value float64
	}
	StackInuse struct {
		Name  string
		Value float64
	}
	StackSys struct {
		Name  string
		Value float64
	}
	Sys struct {
		Name  string
		Value float64
	}
	TotalAlloc struct {
		Name  string
		Value float64
	}
	RandomValue struct {
		Name  string
		Value float64
	}
}
type Counter struct {
	PollCount struct {
		Name  string
		Value int64
	}
}

var GaugeData = map[string]float64{}
var CounterData = map[string]int64{}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

var MetricsStore = map[string]Metrics{}
