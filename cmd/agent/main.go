package main

import (
	"bytes"
	"fmt"
	"github.com/CyrilSbrodov/metricService.git/internal/handlers"
	"io"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"time"
)

type Arg struct {
	pollInterval   time.Duration
	reportInterval time.Duration
}

func main() {
	var arg = Arg{
		2 * time.Second,
		10 * time.Second,
	}
	var gauge handlers.Gauge
	var counter handlers.Counter
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
	ticker := time.NewTicker(arg.reportInterval)
	client := &http.Client{}

	for {
		select {
		case <-ticker.C:
			//сбор метрики
			runtime.ReadMemStats(&memory)
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

			//сбор данных в строку
			val := reflect.ValueOf(gauge)
			result := ""
			for i := 0; i < val.NumField(); i++ {
				result += fmt.Sprintf("%s: %f\n", val.Field(i).Field(0).Interface().(string), val.Field(i).Field(1).Interface().(float64))
			}
			result += fmt.Sprintf("%s: %d", counter.PollCount.Name, counter.PollCount.Value)

			//отправка данных по адресу
			request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/update", bytes.NewBufferString(result))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			//чтение ответа
			response, err := client.Do(request)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(string(body))
			response.Body.Close()

		}
	}

}
