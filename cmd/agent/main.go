package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	//"flag"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

func main() {
	//flagAddress, flagPollInterval, flagReportInterval, flagHash := config.AgentFlagsInit()
	//flag.Parse()
	//
	//cfg := config.NewConfigAgent(*flagAddress, *flagPollInterval, *flagReportInterval, *flagHash)

	cfg := config.AgentConfigInit()

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
			metricsStore = update(metricsStore, count, &cfg)
		}
	}
}

//сбор метрики
func update(store map[string]storage.Metrics, count int64, cfg *config.AgentConfig) map[string]storage.Metrics {
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

		if m.Value == nil {
			m.Value = &value
		}

		m.Hash = hashing(cfg, &m)
		store[m.ID] = m

	}
	var m = storage.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &count,
	}

	m.Hash = hashing(cfg, &m)
	store[m.ID] = m

	value := rand.Intn(256)
	v := float64(value)
	var randomValue = storage.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Value: &v,
	}
	randomValue.Hash = hashing(cfg, &randomValue)
	store[randomValue.ID] = randomValue

	return store
}

//отправка метрики
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
			//fmt.Println(err)
			break
		}
		//_, err = ioutil.ReadAll(resp.Body)
		//if err != nil {
		//	fmt.Println(err)
		//	break
		//}
		resp.Body.Close()
	}
}

func hashing(cfg *config.AgentConfig, m *storage.Metrics) string {
	var hash string
	switch m.MType {
	case "counter":
		hash = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
	case "gauge":
		hash = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
	}
	h := hmac.New(sha256.New, []byte(cfg.Hash))
	h.Write([]byte(hash))
	return fmt.Sprintf("%x", h.Sum(nil))
}
