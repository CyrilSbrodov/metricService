package main

import (
	"fmt"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

type Arg struct {
	pollInterval   time.Duration
	reportInterval time.Duration
}

func main() {
	urlGauge := "http://localhost:8080/update/gauge/"
	urlCounter := "http://localhost:8080/update/counter/"
	var arg = Arg{
		2 * time.Second,
		3 * time.Second,
	}
	var gauge storage.Gauge
	var counter storage.Counter
	var memory runtime.MemStats

	//присвоение имени метрики
	gauge.Alloc.Name = "Alloc"
	gauge.BuckHashSys.Name = "BuckHashSys"
	gauge.Frees.Name = "Frees"
	gauge.GCCPUFraction.Name = "GCCPUFraction"
	gauge.GCSys.Name = "GCSys"
	gauge.HeapAlloc.Name = "HeapAlloc"
	gauge.HeapIdle.Name = "HeapIdle"
	gauge.HeapInuse.Name = "HeapInuse"
	gauge.HeapObjects.Name = "HeapObjects"
	gauge.HeapReleased.Name = "HeapReleased"
	gauge.HeapSys.Name = "HeapSys"
	gauge.LastGC.Name = "LastGC"
	gauge.Lookups.Name = "Lookups"
	gauge.MCacheInuse.Name = "MCacheInuse"
	gauge.MCacheSys.Name = "MCacheSys"
	gauge.MSpanInuse.Name = "MSpanInuse"
	gauge.MSpanSys.Name = "MSpanSys"
	gauge.Mallocs.Name = "Mallocs"
	gauge.NextGC.Name = "NextGC"
	gauge.NumForcedGC.Name = "NumForcedGC"
	gauge.NumGC.Name = "NumGC"
	gauge.OtherSys.Name = "OtherSys"
	gauge.PauseTotalNs.Name = "PauseTotalNs"
	gauge.StackInuse.Name = "StackInuse"
	gauge.StackSys.Name = "StackSys"
	gauge.Sys.Name = "Sys"
	gauge.TotalAlloc.Name = "TotalAlloc"
	counter.PollCount.Name = "PollCount"
	gauge.RandomValue.Name = "RandomValue"

	//запуск тикера
	tickerUpload := time.NewTicker(arg.reportInterval)
	tickerUpdate := time.NewTicker(arg.pollInterval)
	client := &http.Client{}

	for {
		select {
		//отправка метрики 10 сек
		case <-tickerUpload.C:

			val := reflect.ValueOf(gauge)
			for i := 0; i < val.NumField(); i++ {
				//отправка данных по адресу
				value := fmt.Sprintf("%f", val.Field(i).Field(1).Interface().(float64))
				uploadGauge(client, getURL(urlGauge, val.Field(i).Field(0).Interface().(string), value))
			}
			//отправка данных по адресу
			value := strconv.FormatInt(counter.PollCount.Value, 32)

			uploadCounter(client, getURL(urlCounter, counter.PollCount.Name, value))

			//обновление метрики 2 сек
		case <-tickerUpdate.C:
			gauge, counter = update(&memory, gauge, counter)
			fmt.Println(counter)
		}
	}

}

func update(memory *runtime.MemStats, gauge storage.Gauge, counter storage.Counter) (storage.Gauge, storage.Counter) {
	//сбор метрики
	runtime.ReadMemStats(memory)
	gauge.Alloc.Value = float64(memory.Alloc)
	gauge.BuckHashSys.Value = float64(memory.BuckHashSys)
	gauge.Frees.Value = float64(memory.Frees)
	gauge.GCCPUFraction.Value = memory.GCCPUFraction
	gauge.GCSys.Value = float64(memory.GCSys)
	gauge.HeapAlloc.Value = float64(memory.HeapAlloc)
	gauge.HeapIdle.Value = float64(memory.HeapIdle)
	gauge.HeapInuse.Value = float64(memory.HeapInuse)
	gauge.HeapObjects.Value = float64(memory.HeapObjects)
	gauge.HeapReleased.Value = float64(memory.HeapReleased)
	gauge.HeapSys.Value = float64(memory.HeapSys)
	gauge.LastGC.Value = float64(memory.LastGC)
	gauge.Lookups.Value = float64(memory.Lookups)
	gauge.MCacheInuse.Value = float64(memory.MCacheInuse)
	gauge.MCacheSys.Value = float64(memory.MCacheSys)
	gauge.MSpanInuse.Value = float64(memory.MSpanInuse)
	gauge.MSpanSys.Value = float64(memory.MSpanSys)
	gauge.Mallocs.Value = float64(memory.Mallocs)
	gauge.NextGC.Value = float64(memory.NextGC)
	gauge.NumForcedGC.Value = float64(memory.NumForcedGC)
	gauge.NumGC.Value = float64(memory.NumGC)
	gauge.OtherSys.Value = float64(memory.OtherSys)
	gauge.PauseTotalNs.Value = float64(memory.PauseTotalNs)
	gauge.StackInuse.Value = float64(memory.StackInuse)
	gauge.StackSys.Value = float64(memory.StackSys)
	gauge.Sys.Value = float64(memory.Sys)
	gauge.TotalAlloc.Value = float64(memory.TotalAlloc)
	counter.PollCount.Value += 1
	gauge.RandomValue.Value = 1

	return gauge, counter
}

func getURL(url, name, value string) string {
	url += name + "/" + value
	return url
}

func uploadGauge(client *http.Client, url string) {
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	request.Header.Add("Content-Type", "text/plain")

	//чтение ответа
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//fmt.Println(response.Status)

	defer response.Body.Close()
}

func uploadCounter(client *http.Client, url string) {

	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	request.Header.Add("Content-Type", "text/plain")

	//чтение ответа
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//fmt.Println(response.Status)

	defer response.Body.Close()
}
