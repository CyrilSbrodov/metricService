package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

var (
	flagAddress        *string
	flagPollInterval   *string
	flagReportInterval *string
)

func init() {
	flagAddress = flag.String("a", "localhost:8080", "address of service")
	flagPollInterval = flag.String("p", "2s", "update interval")
	flagReportInterval = flag.String("r", "10s", "upload interval to server")

}

func main() {
	flag.Parse()
	cfg := config.NewConfigAgent(*flagAddress, *flagPollInterval, *flagReportInterval)
	fmt.Println(cfg)

	client := &http.Client{}

	var count int64
	metricsStore := storage.MetricsStore

	//запуск тикера
	tickerUpload := time.NewTicker(cfg.ReportInterval)
	tickerUpdate := time.NewTicker(cfg.PollInterval)

	for {
		select {
		//отправка метрики 10 сек
		case <-tickerUpload.C:
			//отправка данных по адресу
			upload(client, cfg.Addr, metricsStore)
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

	value := rand.Intn(256)
	v := float64(value)
	var randomValue = storage.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Value: &v,
	}
	store[randomValue.ID] = randomValue

	return store
}

func upload(client *http.Client, url string, store map[string]storage.Metrics) {

	for _, m := range store {
		metricsJSON, errJSON := json.Marshal(m)
		if errJSON != nil {
			fmt.Println(errJSON)
			break
		}
		req, err := http.NewRequest(http.MethodPost, "http://"+url+"/update/", bytes.NewBuffer(metricsJSON))

		if err != nil {
			fmt.Println(err)
			break
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			break
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			break
		}
		resp.Body.Close()
	}
}

//req.Close = true
//req.Header.Set("Content-Type", "application/json")
//req.Header.Add("Accept", "application/json")

//fmt.Println("перед ду")
//resp, err := client.Do(req)
//if err != nil {
//	fmt.Println(err)
//	os.Exit(1)
//}
