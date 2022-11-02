package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

type Handlers interface {
	Register(router *chi.Mux)
}

type Handler struct {
	storage.Storage
}

// создание роутеров
func (h *Handler) Register(r *chi.Mux) {
	r.Post("/value/", gzipHandle(h.GetHandlerJSON()))
	r.Get("/value/*", gzipHandle(h.GetHandler()))
	r.Get("/", gzipHandle(h.GetAllHandler()))
	r.Post("/update/", gzipHandle(h.CollectHandler()))
	r.Post("/update/gauge/*", gzipHandle(h.GaugeHandler()))
	r.Post("/update/counter/*", gzipHandle(h.CounterHandler()))
	r.Post("/*", gzipHandle(h.OtherHandler()))
	r.Get("/ping", h.PingDB())
}

func NewHandler(storage storage.Storage) Handlers {
	return &Handler{
		storage,
		//*db,
	}
}

//db *repositories.DB

//хендлер получения метрик
func (h *Handler) CollectHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		var m storage.Metrics
		if err := json.Unmarshal(content, &m); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		err = h.Storage.CollectMetrics(m)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Println(err)
			rw.Write([]byte(err.Error()))
			return
		}

		metric, err := h.Storage.GetMetric(m)

		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(err.Error()))
			return
		}
		mJSON, errJSON := json.Marshal(metric)
		if errJSON != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(errJSON.Error()))
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(mJSON)
	}
}

//хендлер получения всех данных
func (h *Handler) GetAllHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		result := h.GetAll()

		rw.Header().Set("Content-Type", "text/html")

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(result))
	}
}

//хендлер получения метрики Gauge
func (h *Handler) GaugeHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		//проверка и разбивка URL
		url := strings.Split(r.URL.Path, "/")
		if len(url) < 5 {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("not value"))
			return
		}

		method := url[1]
		if method != "update" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("incorrect method"))
			return
		}
		types := url[2]

		if types != "gauge" {
			rw.WriteHeader(http.StatusNotImplemented)
			rw.Write([]byte("incorrect type"))
			return
		}
		name := url[3]
		value, err := strconv.ParseFloat(url[4], 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("incorrect value"))
			return
		}

		//отправка значений в БД
		err = h.CollectOrChangeGauge(name, value)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.WriteHeader(http.StatusOK)
	}
}

//хендлер получения метрики Counter
func (h *Handler) CounterHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		//проверка и разбивка URL
		url := strings.Split(r.URL.Path, "/")

		if len(url) < 5 {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("not value"))
			return
		}
		method := url[1]
		if method != "update" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("not value"))
			return
		}
		types := url[2]
		if types != "counter" {
			rw.WriteHeader(http.StatusNotImplemented)
			rw.Write([]byte("incorrect type"))
			return
		}
		name := url[3]

		value, err := strconv.Atoi(url[4])
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}

		//отправка значений в БД
		err = h.CollectOrIncreaseCounter(name, int64(value))
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.WriteHeader(http.StatusOK)

	}
}

//проверка на правильность заполнения update and gauge and counter
func (h *Handler) OtherHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		//проверка и разбивка URL
		url := strings.Split(r.URL.Path, "/")

		if len(url) < 3 {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("not value"))
			return
		}
		method := url[1]
		if method != "update" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("method is wrong"))
			return
		}
		types := url[2]
		if types != "counter" {
			rw.WriteHeader(http.StatusNotImplemented)
			rw.Write([]byte("incorrect type"))
			return
		} else if types != "gauge" {
			rw.WriteHeader(http.StatusNotImplemented)
			rw.Write([]byte("incorrect type"))
			return
		}

		rw.WriteHeader(http.StatusBadRequest)
	}
}

//хендлер получения данных из gauge and counter в формате JSON
func (h *Handler) GetHandlerJSON() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		var m storage.Metrics
		if err := json.Unmarshal(content, &m); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		m, err = h.Storage.GetMetric(m)
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)

		//отправка обновленных метрик обратно
		mJSON, errJSON := json.Marshal(m)
		if errJSON != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(errJSON.Error()))
			return
		}
		rw.Write(mJSON)
	}
}

//хендлер получения данных из gauge and counter
func (h *Handler) GetHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		//проверка и разбивка URL
		url := strings.Split(r.URL.Path, "/")
		if len(url) < 3 {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("incorrect request"))
			return
		}
		method := url[1]
		if method != "value" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("method is wrong"))
			return
		}
		name := url[3]
		types := url[2]
		if types == "gauge" {

			//получение значений из gauge
			value, err := h.GetGauge(name)
			if err != nil {
				rw.WriteHeader(http.StatusNotFound)
				rw.Write([]byte("incorrect name"))
				return
			}
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(fmt.Sprintf("%v", value)))
			return
		} else if types == "counter" {

			//получение значений из counter
			value, err := h.GetCounter(name)
			if err != nil {
				rw.WriteHeader(http.StatusNotFound)
				rw.Write([]byte("incorrect name"))
				return
			}
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(fmt.Sprintf("%v", value)))
			return
		} else {
			rw.WriteHeader(http.StatusNotImplemented)
			rw.Write([]byte("incorrect type"))
			return
		}
	}
}

//func (h Handler) Pings() http.HandlerFunc {
//	return func(rw http.ResponseWriter, r *http.Request) {
//		fmt.Println("ping")
//		err := h.PingClient(context.Background())
//		if err != nil {
//			rw.WriteHeader(http.StatusInternalServerError)
//			rw.Write([]byte(err.Error()))
//			return
//		}
//		rw.Header().Set("Content-Type", "text/html")
//		rw.WriteHeader(http.StatusOK)
//		rw.Write([]byte("200"))
//	}
//}

func (h *Handler) PingDB() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("ping")
		errPing := h.PingClient()
		if errPing != nil {
			fmt.Println(errPing)
			fmt.Println("error repo")
		}
		err := h.Storage.PingClient()
		if err != nil {
			fmt.Println("error pgs")
			http.Error(rw, "", http.StatusInternalServerError)
		}
		//rw.Header().Set("Content-Type", "text/html")
		rw.WriteHeader(http.StatusOK)
	}
}
