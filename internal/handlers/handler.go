package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
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
	//storage.Service
	storage.Storage
}

// создание роутеров
func (h Handler) Register(r *chi.Mux) {
	r.Post("/value", h.GetHandlerJSON())
	r.Get("/value/*", h.GetHandler())
	r.Get("/", h.GetAllHandler())
	r.Post("/update", h.CollectHandler())
	r.Post("/update/gauge/*", h.GaugeHandler())
	r.Post("/update/counter/*", h.CounterHandler())
	r.Post("/*", h.OtherHandler())

}

func NewHandler(storage storage.Storage) Handlers {
	return &Handler{
		storage,
	}
}

//хендлер получения метрик
func (h Handler) CollectHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		//fmt.Println("CollectHandler")
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		}
		defer r.Body.Close()
		var m storage.Metrics
		if err := json.Unmarshal(content, &m); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		}
		h.Storage.CollectMetrics(m)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
	}
}

//хендлер получения всех данных
func (h Handler) GetAllHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		//fmt.Println("GetAllHandler")
		t, err := template.ParseFiles("index.html")
		if err != nil {
			log.Print("template parsing error: ", err)
		}
		t.Execute(rw, nil)
		result := h.GetAll()
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(result))
	}
}

//хендлер получения метрики Gauge
func (h Handler) GaugeHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		//fmt.Println("GaugeHandler")
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
func (h Handler) CounterHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		//fmt.Println("CounterHandler")
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
func (h Handler) OtherHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		//fmt.Println("OtherHandler")
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

//хендлер получения данных из gauge and counter
func (h Handler) GetHandlerJSON() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		//fmt.Println("GetHandlerJSON")
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		}
		defer r.Body.Close()
		var m storage.Metrics
		if err := json.Unmarshal(content, &m); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		}
		m, err = h.Storage.GetMetric(&m)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		//TODO
		if m.Value != nil {
			rw.Write([]byte(fmt.Sprintf("%s : %f", m.ID, *m.Value)))
		} else if m.Delta != nil {
			rw.Write([]byte(fmt.Sprintf("%s : %d", m.ID, *m.Delta)))
		} else {
			rw.Write([]byte(fmt.Sprintf("%s : %d, %d", m.ID, m.Value, m.Delta)))
		}
	}
}

//хендлер получения данных из gauge and counter
func (h Handler) GetHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		//fmt.Println("GetHandler")
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
