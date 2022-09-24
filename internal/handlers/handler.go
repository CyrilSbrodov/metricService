package handlers

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
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
