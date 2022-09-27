package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
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

	r.Post("/update/gauge/*", h.GaugeHandler())
	r.Post("/update/counter/*", h.CounterHandler())
	r.Post("/*", h.OtherHandler())
	r.Get("/value/*", h.GetHandler())
	r.Get("/", h.GetAllHandler())

}

func NewHandler(storage storage.Storage) Handlers {
	return &Handler{
		storage,
	}
}

//хендлер из задания про родственников
func (h Handler) UserViewHandler(users map[string]storage.User) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if r.URL.Query().Get("user_id") == "" {
			http.Error(rw, "userId is empty", http.StatusBadRequest)
			return
		}

		user, ok := users[userID]
		if !ok {
			http.Error(rw, "user not found", http.StatusNotFound)
			return
		}

		jsonUser, err := json.Marshal(user)
		if err != nil {
			http.Error(rw, "can't provide a json. internal error", http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(jsonUser)
	}
}

//хендлер получения метрики Gauge
func (h Handler) GaugeHandler() http.HandlerFunc {
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
		err = h.CollectGauge(name, value)
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
func (h Handler) GetHandler() http.HandlerFunc {
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

//хендлер получения всех данных из gauge and counter
func (h Handler) GetAllHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
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
