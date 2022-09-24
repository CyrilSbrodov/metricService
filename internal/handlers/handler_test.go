package handlers_test

import (
	"encoding/json"
	"fmt"
	"github.com/CyrilSbrodov/metricService.git/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserViewHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		user        handlers.User
	}

	tests := []struct {
		name    string
		request string
		users   map[string]handlers.User
		want    want
	}{
		{
			name: "simple test #1 'ok'",
			users: map[string]handlers.User{
				"1": {
					ID:        "1",
					FirstName: "Misha",
					LastName:  "Popov",
				},
			},
			want: want{
				contentType: "application/json",
				statusCode:  200,
				user: handlers.User{
					"1",
					"Misha",
					"Popov",
				},
			},
			request: "/users?user_id=1",
		},
		{
			name: "simple test #2 'not found'",
			users: map[string]handlers.User{
				"1": {
					ID:        "1",
					FirstName: "Misha",
					LastName:  "Popov",
				},
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
				user:        handlers.User{},
			},
			request: "/users?user_id=3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handlers.UserViewHandler(tt.users))
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			//fmt.Println(result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			var user handlers.User
			//fmt.Println(handlers.User{})
			err = json.Unmarshal(userResult, &user)
			fmt.Println("unmarshal", user)
			//require.NoError(t, err)

			assert.Equal(t, tt.want.user, user)
		})
	}
}

//type userbad struct {
//	users map[string]handlers.User
//}
//
//func (u *userbad) GetID(id string) error {
//	return errors.New("Bad connection")
//}
//
//func TestUserViewHandler500(t *testing.T) {
//
//	type want struct {
//		contentType string
//		statusCode  int
//		user        handlers.User
//	}
//
//	tests := []struct {
//		name    string
//		request string
//		users   userbad
//		want    want
//	}{
//		{
//			name: "simple test #3 'bad request'",
//			users: userbad{
//				map[string]handlers.User{
//					"1": {
//						"1",
//						"Misha",
//						"Ivanov",
//					},
//				},
//			},
//			want: want{
//				contentType: "text/plain; charset=utf-8",
//				statusCode:  500,
//				user:        handlers.User{},
//			},
//			request: "/users?user_id=1",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
//			w := httptest.NewRecorder()
//			h := http.HandlerFunc(handlers.UserViewHandler(tt.users.users))
//			h.ServeHTTP(w, request)
//			result := w.Result()
//
//			assert.Equal(t, tt.want.statusCode, result.StatusCode)
//			fmt.Println(result.StatusCode)
//			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
//
//			userResult, err := ioutil.ReadAll(result.Body)
//			require.NoError(t, err)
//			err = result.Body.Close()
//			require.NoError(t, err)
//
//			var user handlers.User
//			//fmt.Println(handlers.User{})
//			err = json.Unmarshal(userResult, &user)
//			fmt.Println("unmarshal", user)
//			require.NoError(t, err)
//
//			assert.Equal(t, tt.want.user, user)
//		})
//	}
//}
