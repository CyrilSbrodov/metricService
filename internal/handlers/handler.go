package handlers

import (
	"encoding/json"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"io/ioutil"
	"net/http"
	"strings"
)

type Handlers interface {
	Register(router *http.ServeMux)
}

type handler struct {
	//storage.Service
	storage.Storage
}

func (h handler) Register(router *http.ServeMux) {
	//router.HandleFunc("/user", h.UserViewHandler(users map[string]storage.User))
	router.HandleFunc("/update/", h.UpdateHandler())
}

func NewHandler(storage storage.Storage) Handlers {
	return &handler{
		storage,
	}
}

func (h handler) UserViewHandler(users map[string]storage.User) http.HandlerFunc {
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

func (h handler) UpdateHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		defer r.Body.Close()

		url := strings.Split(r.URL.Path, "/")

		err = h.Collect(url[2], url[3], string(content))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.WriteHeader(http.StatusOK)
		//fmt.Println(string(content))

	}
}
