package main

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/cmd/loggers"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

func main() {
	wg := &sync.WaitGroup{}
	cfg := config.AgentConfigInit()
	client := &http.Client{}

	var count int64
	metrics := storage.NewAgentMetrics()
	logger := loggers.NewLogger()
	//запуск тикера
	tickerUpload := time.NewTicker(cfg.ReportInterval)
	tickerUpdate := time.NewTicker(cfg.PollInterval)

	for {
		select {
		//отправка метрики 10 сек
		case <-tickerUpload.C:
			//отправка данных по адресу
			wg.Add(2)
			go upload(client, cfg.Addr, metrics, logger, wg)
			go uploadBatch(client, cfg.Addr, metrics, logger, wg)
			//обновление метрики 2 сек
		case <-tickerUpdate.C:
			count++
			wg.Add(2)
			go update(metrics, count, &cfg, wg)
			go updateOtherMetrics(metrics, wg, &cfg)
		}
	}
	wg.Wait()
}

//сбор метрики
func update(store *storage.AgentMetrics, count int64, cfg *config.AgentConfig, wg *sync.WaitGroup) {
	store.Sync.Lock()
	defer store.Sync.Unlock()
	defer wg.Done()

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
		store.Store[m.ID] = m

	}
	var m = storage.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &count,
	}

	m.Hash = hashing(cfg, &m)
	store.Store[m.ID] = m

	value := rand.Intn(256)
	v := float64(value)
	var randomValue = storage.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Value: &v,
	}
	randomValue.Hash = hashing(cfg, &randomValue)
	store.Store[randomValue.ID] = randomValue
}

//отправка метрики
func upload(client *http.Client, url string, store *storage.AgentMetrics, logger *loggers.Logger, wg *sync.WaitGroup) {
	store.Sync.Lock()
	defer store.Sync.Unlock()
	defer wg.Done()
	for _, m := range store.Store {
		metricsJSON, errJSON := json.Marshal(m)
		if errJSON != nil {
			fmt.Println(errJSON)
			break
		}
		req, err := http.NewRequest(http.MethodPost, "http://"+url+"/update/", bytes.NewBuffer(metricsJSON))

		if err != nil {
			logger.LogErr(err, "Failed to request")
			fmt.Println(err)
			break
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			logger.LogErr(err, "Failed to do request")
			break
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.LogErr(err, "Failed to read body")
			break
		}
		resp.Body.Close()
	}
}

func uploadBatch(client *http.Client, url string, store *storage.AgentMetrics, logger *loggers.Logger, wg *sync.WaitGroup) {
	store.Sync.Lock()
	defer store.Sync.Unlock()
	defer wg.Done()
	var metrics []storage.Metrics
	for _, m := range store.Store {
		metrics = append(metrics, m)
	}
	metricsJSON, errJSON := json.Marshal(metrics)
	if errJSON != nil {
		logger.LogErr(errJSON, "Failed to Marshal metrics to JSON")
		fmt.Println(errJSON)
		return
	}
	metricsCompress, err := compress(metricsJSON, logger)
	if err != nil {
		logger.LogErr(errJSON, "Failed to compress metrics")
		return
	}
	if len(metricsCompress) == 0 {
		return
	}
	req, err := http.NewRequest(http.MethodPost, "http://"+url+"/updates/", bytes.NewBuffer(metricsCompress))

	if err != nil {
		logger.LogErr(err, "Failed to request")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	resp, err := client.Do(req)
	if err != nil {
		logger.LogErr(err, "Failed to do request")
		return
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogErr(err, "Failed to read body")
		return
	}
	resp.Body.Close()
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

func compress(store []byte, logger *loggers.Logger) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	_, err := w.Write(store)
	if err != nil {
		logger.LogErr(err, "failed to write data to compress temporary buffer")
		return nil, fmt.Errorf("failed to write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		logger.LogErr(err, "failed compress data")
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	return b.Bytes(), nil
}

func updateOtherMetrics(store *storage.AgentMetrics, wg *sync.WaitGroup, cfg *config.AgentConfig) {
	store.Sync.Lock()
	defer store.Sync.Unlock()
	defer wg.Done()
	var value float64
	v, _ := mem.VirtualMemory()
	cpuValue, _ := cpu.Percent(0, false)
	value = float64(v.Total)
	var total = storage.Metrics{
		ID:    "TotalMemory",
		MType: "gauge",
		Value: &value,
	}
	value = float64(v.Free)
	var free = storage.Metrics{
		ID:    "FreeMemory",
		MType: "gauge",
		Value: &value,
	}
	var cpu = storage.Metrics{
		ID:    "CPUutilization1",
		MType: "gauge",
		Value: &cpuValue[0],
	}
	total.Hash = hashing(cfg, &total)
	free.Hash = hashing(cfg, &free)
	cpu.Hash = hashing(cfg, &cpu)
	store.Store[total.ID] = total
	store.Store[free.ID] = free
	store.Store[cpu.ID] = cpu
}
