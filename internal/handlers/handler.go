package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
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

	r.Post("/*", h.OtherHandler())
	r.Post("/value", h.GetHandler())
	r.Get("/", h.GetAllHandler())
	r.Post("/update", h.CollectHandler())

}

func NewHandler(storage storage.Storage) Handlers {
	return &Handler{
		storage,
	}
}

//хендлер получения метрик
func (h Handler) CollectHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		}
		defer r.Body.Close()
		var store map[string]storage.Metrics
		if err := json.Unmarshal(content, &store); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		}
		h.Storage.CollectMetrics(store)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
	}
}

//хендлер получения данных
func (h Handler) GetHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
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

//хендлер получения всех данных
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

//проверка на правильность заполнения url
func (h Handler) OtherHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		//проверка и разбивка URL
		url := strings.Split(r.URL.Path, "/")
		method := url[1]
		if method != "update" || method != "value" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("method is wrong"))
			return
		}
		rw.WriteHeader(http.StatusBadRequest)
	}
}
