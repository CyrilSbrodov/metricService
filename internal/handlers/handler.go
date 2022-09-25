package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
}

//type Gauge float64
//type Counter int64

type Gauge struct {
	Alloc struct {
		Name  string
		Value float64
	}
	BuckHashSys struct {
		Name  string
		Value float64
	}
	Frees struct {
		Name  string
		Value float64
	}
	GCCPUFraction struct {
		Name  string
		Value float64
	}
	GCSys struct {
		Name  string
		Value float64
	}
	HeapAlloc struct {
		Name  string
		Value float64
	}
	HeapIdle struct {
		Name  string
		Value float64
	}
	HeapInuse struct {
		Name  string
		Value float64
	}
	HeapObjects struct {
		Name  string
		Value float64
	}
	HeapReleased struct {
		Name  string
		Value float64
	}
	HeapSys struct {
		Name  string
		Value float64
	}
	LastGC struct {
		Name  string
		Value float64
	}
	Lookups struct {
		Name  string
		Value float64
	}
	MCacheInuse struct {
		Name  string
		Value float64
	}
	MCacheSys struct {
		Name  string
		Value float64
	}
	MSpanInuse struct {
		Name  string
		Value float64
	}
	MSpanSys struct {
		Name  string
		Value float64
	}
	Mallocs struct {
		Name  string
		Value float64
	}
	NextGC struct {
		Name  string
		Value float64
	}
	NumForcedGC struct {
		Name  string
		Value float64
	}
	NumGC struct {
		Name  string
		Value float64
	}
	OtherSys struct {
		Name  string
		Value float64
	}
	PauseTotalNs struct {
		Name  string
		Value float64
	}
	StackInuse struct {
		Name  string
		Value float64
	}
	StackSys struct {
		Name  string
		Value float64
	}
	Sys struct {
		Name  string
		Value float64
	}
	TotalAlloc struct {
		Name  string
		Value float64
	}
	RandomValue struct {
		Name  string
		Value float64
	}
}
type Counter struct {
	PollCount struct {
		Name  string
		Value int64
	}
}

func UserViewHandler(users map[string]User) http.HandlerFunc {
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

func CollectMetricData(g *Gauge, c *Counter) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		str := fmt.Sprintf("%f "+"%d ", g.Sys, c.PollCount)
		resultJson, err := json.MarshalIndent(str, " ", " ")
		if err != nil {
			errors.New(fmt.Sprintf("не удалось перекодировать данные. ошибка: %v", err))
		}
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(resultJson)
	}
}

func UpdateHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		//maps := make(map[string]string)
		defer r.Body.Close()
		//stringss := strings.Split(string(content), "\n")
		//for _, s := range stringss {
		//	ss := strings.Split(s, ":")
		//	maps[ss[0]] = ss[1]
		//}
		fmt.Println(string(content))

	}
}
