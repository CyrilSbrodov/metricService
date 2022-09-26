package storage

type User struct {
	ID        string
	FirstName string
	LastName  string
}

//type Gauge float64
//type Counter int64

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

var GaugeData = map[string]float64{
	"Alloc":         0,
	"BuckHashSys":   0,
	"Frees":         0,
	"GCCPUFraction": 0,
	"GCSys":         0,
	"HeapAlloc":     0,
	"HeapIdle":      0,
	"HeapInuse":     0,
	"HeapObjects":   0,
	"HeapReleased":  0,
	"HeapSys":       0,
	"LastGC":        0,
	"Lookups":       0,
	"MCacheInuse":   0,
	"MCacheSys":     0,
	"MSpanInuse":    0,
	"MSpanSys":      0,
	"Mallocs":       0,
	"NextGC":        0,
	"NumForcedGC":   0,
	"NumGC":         0,
	"OtherSys":      0,
	"PauseTotalNs":  0,
	"StackInuse":    0,
	"StackSys":      0,
	"Sys":           0,
	"TotalAlloc":    0,
	"RandomValue":   0,
}

var CounterData = map[string]int64{
	"PollCount": 0,
}
