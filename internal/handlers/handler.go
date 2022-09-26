package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Handlers interface {
	Register(router *chi.Mux)
}

type Handler struct {
	//storage.Service
	storage.Storage
}

func (h Handler) Register(r *chi.Mux) {
	//router.HandleFunc("/user", h.UserViewHandler(users map[string]storage.User))
	r.Group(func(r chi.Router) {
		r.Post("/update/gauge/*", h.GaugeHandler())
		r.Post("/update/counter/*", h.CounterHandler())
		r.Post("/*", h.OtherHandler())
		r.Get("/value/*", h.GetHandler())
	})
}

func NewHandler(storage storage.Storage) Handlers {
	return &Handler{
		storage,
	}
}

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

//func CollectMetricData(g *Gauge, c *Counter) http.HandlerFunc {
//	return func(rw http.ResponseWriter, r *http.Request) {
//		str := fmt.Sprintf("%f "+"%d ", g.Sys, c.PollCount)
//		resultJson, err := json.MarshalIndent(str, " ", " ")
//		if err != nil {
//			errors.New(fmt.Sprintf("не удалось перекодировать данные. ошибка: %v", err))
//		}
//		rw.Header().Set("Access-Control-Allow-Origin", "*")
//		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
//		rw.WriteHeader(http.StatusOK)
//		_, _ = rw.Write(resultJson)
//	}
//}

func (h Handler) GaugeHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		defer r.Body.Close()

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
		err = h.CollectGauge(name, value)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.WriteHeader(http.StatusOK)
	}
}

func (h Handler) CounterHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		defer r.Body.Close()

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

		err = h.CollectCounter(name, int64(value))
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.WriteHeader(http.StatusOK)
	}
}
func (h Handler) OtherHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		defer r.Body.Close()

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

func (h Handler) GetHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
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
