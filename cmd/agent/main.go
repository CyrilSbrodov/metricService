package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

type Arg struct {
	pollInterval   time.Duration
	reportInterval time.Duration
}

func main() {
	url := "http://localhost:8080/update"
	var arg = Arg{
		2 * time.Second,
		10 * time.Second,
	}
	var count int64
	metricsStore := storage.MetricsStore

	//запуск тикера
	tickerUpload := time.NewTicker(arg.reportInterval)
	tickerUpdate := time.NewTicker(arg.pollInterval)
	client := &http.Client{}

	for {
		select {
		//отправка метрики 10 сек
		case <-tickerUpload.C:
			//отправка данных по адресу
			upload(client, url, metricsStore)
			//обновление метрики 2 сек
		case <-tickerUpdate.C:
			count++
			metricsStore = update(metricsStore, count)
		}
	}
}

func update(store map[string]storage.Metrics, count int64) map[string]storage.Metrics {
	//сбор метрики
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	val := reflect.ValueOf(memory)
	for i := 0; i < val.NumField(); i++ {
		var value float64
		var m storage.Metrics
		m.ID = reflect.TypeOf(memory).Field(i).Name
		m.MType = "gauge"

		metricValue := val.Field(i).Interface()
		switch valueType := metricValue.(type) {
		case float64:
			value = valueType
			m.Value = &value
		case uint64:
			var x uint64
			x = valueType
			value = float64(x)
			m.Value = &value
		case uint32:
			var x uint32
			x = valueType
			value = float64(x)
			m.Value = &value
		}
		store[m.ID] = m

	}
	var m = storage.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &count,
	}
	store[m.ID] = m
	return store
}

func upload(client *http.Client, url string, store map[string]storage.Metrics) {

	metricsJSON, err := json.Marshal(store)
	if err != nil {
		fmt.Errorf("не удалось перекодировать данные. ошибка: %v", err)
	}

	req, err := http.Post(url, "application/json", bytes.NewBuffer(metricsJSON))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer req.Body.Close()

}
