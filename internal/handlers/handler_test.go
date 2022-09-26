package handlers_test

import (
	"github.com/CyrilSbrodov/metricService.git/internal/handlers"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"github.com/CyrilSbrodov/metricService.git/internal/storage/repositories"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func Test_handler_UserViewHandler(t *testing.T) {
//
//	type want struct {
//		contentType string
//		statusCode  int
//		user        storage.User
//	}
//
//	//type fields struct {
//	//	Storage storage.Storage
//	//}
//
//	tests := []struct {
//		name string
//		request string
//		users map[string]storage.User
//		want  want
//	}{
//		{
//			name: "simple test #1 'ok'",
//			users: map[string]storage.User{
//				"1": {
//					ID:        "1",
//					FirstName: "Misha",
//					LastName:  "Popov",
//				},
//			},
//			want: want{
//				contentType: "application/json",
//				statusCode:  200,
//				user: storage.User{
//					"1",
//					"Misha",
//					"Popov",
//				},
//			},
//			request: "/users?user_id=1",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
//			w := httptest.NewRecorder()
//			h := http.HandlerFunc(handlers.Handler{storage.Storage(tt.users)}(tt.users))
//			h.ServeHTTP(w, request)
//			result := w.Result()
//
//			assert.Equal(t, tt.want.statusCode, result.StatusCode)
//			//fmt.Println(result.StatusCode)
//			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
//
//			userResult, err := ioutil.ReadAll(result.Body)
//			require.NoError(t, err)
//			err = result.Body.Close()
//			require.NoError(t, err)
//
//			var user storage.User
//			//fmt.Println(handlers.User{})
//			err = json.Unmarshal(userResult, &user)
//			//fmt.Println("unmarshal", user)
//			require.NoError(t, err)
//
//			assert.Equal(t, tt.want.user, user)
//		})
//	}
//}

func TestHandler_GaugeHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	repo := repositories.NewRepository()

	type fields struct {
		Storage storage.Storage
	}
	tests := []struct {
		name    string
		request string
		fields  fields
		want    want
	}{
		{
			name: "Test ok",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/gauge/test/100",
			want: want{
				200,
			},
		},
		{
			name: "Test 501",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/gaug/test/100",
			want: want{
				501,
			},
		},
		{
			name: "Test 400",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/gauge/test/none",
			want: want{
				400,
			},
		},
		{
			name: "Test 404",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/gauge/test",
			want: want{
				404,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := handlers.Handler{
				Storage: tt.fields.Storage,
			}
			h.GaugeHandler().ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_CounterHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	repo := repositories.NewRepository()

	type fields struct {
		Storage storage.Storage
	}
	tests := []struct {
		name    string
		request string
		fields  fields
		want    want
	}{
		{
			name: "Test ok",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/counter/test/100",
			want: want{
				200,
			},
		},
		{
			name: "Test 501",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/test/test/100",
			want: want{
				501,
			},
		},
		{
			name: "Test 400",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/counter/test/none",
			want: want{
				400,
			},
		},
		{
			name: "Test 404",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/counter/test",
			want: want{
				404,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := handlers.Handler{
				Storage: tt.fields.Storage,
			}
			h.CounterHandler().ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_OtherHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	repo := repositories.NewRepository()

	type fields struct {
		Storage storage.Storage
	}
	tests := []struct {
		name    string
		request string
		fields  fields
		want    want
	}{
		{
			name: "Test wrong path/method, code 404",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/test/counter/test/100",
			want: want{
				404,
			},
		},
		{
			name: "Test wrong types, code 501",
			fields: fields{
				repo,
			},
			request: "http://localhost:8080/update/test/test/100",
			want: want{
				501,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := handlers.Handler{
				Storage: tt.fields.Storage,
			}
			h.OtherHandler().ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
