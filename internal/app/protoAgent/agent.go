package protoAgent

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/cmd/loggers"
	pb "github.com/CyrilSbrodov/metricService.git/internal/app/proto"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

type Agenter interface {
	NewAgentApp() *AgentApp
}

//AgentApp структура агента
type AgentApp struct {
	closed chan struct{}
	cfg    config.AgentConfig
	logger *loggers.Logger
	wg     sync.WaitGroup
}

//NewAgentApp создание нового агента
func NewAgentApp() *AgentApp {
	cfg := config.AgentConfigInit()
	logger := loggers.NewLogger()
	return &AgentApp{
		closed: make(chan struct{}),
		cfg:    *cfg,
		logger: logger,
	}
}

// Run запуск агента
func (a *AgentApp) Run() {
	conn, err := grpc.Dial(":3202", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		a.logger.LogErr(err, "failed to dial")
	}
	defer conn.Close()

	client := pb.NewStorageClient(conn)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	//запуск тикеров
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.RunTickers(client)
	}()

	//прослушивание канала на сигналы завершения работы
	<-done

	a.logger.LogInfo("", "", "Agent Shutdown")
	//остановка агента, если в канал поступает сигнал
	a.Stop()
	a.logger.LogInfo("", "", "Agent Shutdown gracefully")
}

// RunTickers запуск тикеров
func (a *AgentApp) RunTickers(client pb.StorageClient) {
	var count int64
	metrics := storage.NewAgentMetricsProto()
	//запуск тикера
	tickerUpload := time.NewTicker(a.cfg.ReportInterval)
	tickerUpdate := time.NewTicker(a.cfg.PollInterval)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		//отправка метрики 10 сек
		case <-tickerUpload.C:
			//отправка данных по адресу
			a.wg.Add(2)
			if a.cfg.CryptoPROKey != "" {
				//go a.uploadCrypto(metrics, ctx)
				//go a.uploadBatchCrypto(metrics, ctx)
			} else {
				go a.upload(ctx, metrics, client)
				go a.uploadBatch(ctx, metrics, client)
			}
			//обновление метрики 2 сек
		case <-tickerUpdate.C:
			count++
			a.wg.Add(2)
			go a.update(metrics, count)
			go a.updateOtherMetrics(metrics)
		case <-a.closed:
			return
		}
	}
}

// Stop остановка агента
func (a *AgentApp) Stop() {
	close(a.closed)
	a.wg.Wait()
}

//Сбор метрики
func (a *AgentApp) update(store *storage.AgentMetricsProto, count int64) {
	store.Sync.Lock()
	defer store.Sync.Unlock()
	defer a.wg.Done()

	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	val := reflect.ValueOf(memory)
	for i := 0; i < val.NumField(); i++ {
		var value float64
		var m pb.Metrics
		m.ID = reflect.TypeOf(memory).Field(i).Name
		m.MType = "gauge"

		metricValue := val.Field(i).Interface()
		switch valueType := metricValue.(type) {
		case float64:
			value = valueType
			m.Value = value
		case uint64:
			var x uint64
			x = valueType
			value = float64(x)
			m.Value = value
		case uint32:
			var x uint32
			x = valueType
			value = float64(x)
			m.Value = value
		}

		if m.Value == 0 {
			m.Value = value
		}

		m.Hash = a.hashing(&m)
		store.Store[m.ID] = m

	}
	var m = pb.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: count,
	}

	m.Hash = a.hashing(&m)
	store.Store[m.ID] = m

	value := rand.Intn(256)
	v := float64(value)
	var randomValue = pb.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Value: v,
	}
	randomValue.Hash = a.hashing(&randomValue)
	store.Store[randomValue.ID] = randomValue
}

//Сбор остальных метрик.
func (a *AgentApp) updateOtherMetrics(store *storage.AgentMetricsProto) {
	store.Sync.Lock()
	defer store.Sync.Unlock()
	defer a.wg.Done()
	var value float64
	v, _ := mem.VirtualMemory()
	cpuValue, _ := cpu.Percent(0, false)
	value = float64(v.Total)
	var total = pb.Metrics{
		ID:    "TotalMemory",
		MType: "gauge",
		Value: value,
	}
	value = float64(v.Free)
	var free = pb.Metrics{
		ID:    "FreeMemory",
		MType: "gauge",
		Value: value,
	}
	var cpu = pb.Metrics{
		ID:    "CPUutilization1",
		MType: "gauge",
		Value: cpuValue[0],
	}
	total.Hash = a.hashing(&total)
	free.Hash = a.hashing(&free)
	cpu.Hash = a.hashing(&cpu)
	store.Store[total.ID] = total
	store.Store[free.ID] = free
	store.Store[cpu.ID] = cpu
}

//Отправка метрики
func (a *AgentApp) upload(ctx context.Context, store *storage.AgentMetricsProto, client pb.StorageClient) {
	store.Sync.Lock()
	defer store.Sync.Unlock()
	defer a.wg.Done()
	for _, m := range store.Store {

		resp, err := client.CollectMetric(ctx, a.AddMetric(&m))
		if err != nil {
			a.logger.LogErr(err, "Failed to request")
			fmt.Println(err)
			break
		}
		a.logger.LogInfo("metric:", resp.String(), "")
	}
}

//
////Отправка метрики Crypto
//func (a *AgentApp) uploadCrypto(store *storage.AgentMetrics) {
//	store.Sync.Lock()
//	defer store.Sync.Unlock()
//	defer a.wg.Done()
//	for _, m := range store.Store {
//		metricsJSON, errJSON := json.Marshal(m)
//		if errJSON != nil {
//			fmt.Println(errJSON)
//			break
//		}
//		//шифрование данных с помощью публичного ключа
//		mByte, err := crypto.EncryptedData(metricsJSON, a.public, a.logger)
//		if err != nil {
//			a.logger.LogErr(err, "error from encrypted")
//			break
//		}
//		req, err := http.NewRequest(http.MethodPost, a.url+a.cfg.Addr+"/update/", bytes.NewBuffer(mByte))
//
//		if err != nil {
//			a.logger.LogErr(err, "Failed to request")
//			fmt.Println(err)
//			break
//		}
//		req.Header.Set("X-Real-IP", getIP(req))
//		req.Header.Set("Content-Type", "application/json")
//		req.Header.Add("Accept", "application/json")
//
//		resp, err := a.client.Do(req)
//		if err != nil {
//			a.logger.LogErr(err, "Failed to do request")
//			break
//		}
//		_, err = io.ReadAll(resp.Body)
//		if err != nil {
//			a.logger.LogErr(err, "Failed to read body")
//			break
//		}
//		resp.Body.Close()
//	}
//}

//Отправка метрики батчами.
func (a *AgentApp) uploadBatch(ctx context.Context, store *storage.AgentMetricsProto, client pb.StorageClient) {
	store.Sync.Lock()
	defer store.Sync.Unlock()
	defer a.wg.Done()
	var metrics []*pb.Metrics
	for _, m := range store.Store {
		metrics = append(metrics, &m)
	}
	resp, err := client.CollectMetrics(ctx, a.AddMetrics(metrics))
	if err != nil {
		a.logger.LogErr(err, "Failed to request")
		fmt.Println(err)
		return
	}
	a.logger.LogInfo("metric:", resp.String(), "")
}

////Отправка метрики батчами Crypto.
//func (a *AgentApp) uploadBatchCrypto(store *storage.AgentMetrics) {
//	store.Sync.Lock()
//	defer store.Sync.Unlock()
//	defer a.wg.Done()
//	var metrics []storage.Metrics
//	for _, m := range store.Store {
//		metrics = append(metrics, m)
//	}
//	metricsJSON, errJSON := json.Marshal(metrics)
//	if errJSON != nil {
//		a.logger.LogErr(errJSON, "Failed to Marshal metrics to JSON")
//		fmt.Println(errJSON)
//		return
//	}
//	//шифрование данных с помощью публичного ключа
//	mByte, err := crypto.EncryptedData(metricsJSON, a.public, a.logger)
//	if err != nil {
//		a.logger.LogErr(err, "error from encrypted")
//		return
//	}
//	metricsCompress, err := a.compress(mByte)
//	if err != nil {
//		a.logger.LogErr(errJSON, "Failed to compress metrics")
//		return
//	}
//	if len(metricsCompress) == 0 {
//		return
//	}
//	req, err := http.NewRequest(http.MethodPost, a.url+a.cfg.Addr+"/updates/", bytes.NewBuffer(metricsCompress))
//
//	if err != nil {
//		a.logger.LogErr(err, "Failed to request")
//		return
//	}
//	req.Header.Set("X-Real-IP", getIP(req))
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Add("Content-Encoding", "gzip")
//
//	resp, err := a.client.Do(req)
//	if err != nil {
//		a.logger.LogErr(err, "Failed to do request")
//		return
//	}
//	_, err = io.ReadAll(resp.Body)
//	if err != nil {
//		a.logger.LogErr(err, "Failed to read body")
//		return
//	}
//	resp.Body.Close()
//}

//Функция хеширования.
func (a *AgentApp) hashing(m *pb.Metrics) string {
	var hash string
	switch m.MType {
	case "counter":
		hash = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, m.Delta)
	case "gauge":
		hash = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, m.Value)
	}
	h := hmac.New(sha256.New, []byte(a.cfg.Hash))
	h.Write([]byte(hash))
	return fmt.Sprintf("%x", h.Sum(nil))
}

//Функция сжатия данных.
func (a *AgentApp) compress(store []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	_, err := w.Write(store)
	if err != nil {
		a.logger.LogErr(err, "failed to write data to compress temporary buffer")
		return nil, fmt.Errorf("failed to write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		a.logger.LogErr(err, "failed compress data")
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	return b.Bytes(), nil
}

func getIP(r *http.Request) string {
	IP := r.RemoteAddr
	if IP == "" {
		IP = "127.0.0.1:8080"
	}
	return IP
}

func (a *AgentApp) AddMetric(metrics *pb.Metrics) *pb.AddMetricRequest {
	return &pb.AddMetricRequest{
		Metrics: metrics,
	}
}
func (a *AgentApp) AddMetrics(metrics []*pb.Metrics) *pb.AddMetricsRequest {

	return &pb.AddMetricsRequest{
		Metrics: metrics,
	}
}
