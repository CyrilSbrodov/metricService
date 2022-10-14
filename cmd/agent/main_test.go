package main

//import (
//	"net/http"
//	"reflect"
//	"runtime"
//	"testing"
//
//	"github.com/CyrilSbrodov/metricService.git/internal/storage"
//)
//
//func Test_getURL(t *testing.T) {
//	type args struct {
//		url   string
//		name  string
//		value string
//	}
//	tests := []struct {
//		name string
//		args args
//		want string
//	}{
//		{
//			name: "test ok",
//			args: args{
//				url:   "http://localhost:8080/gauge/",
//				name:  "test",
//				value: "100",
//			},
//			want: "http://localhost:8080/gauge/test/100",
//		},
//		{
//			name: "test not ok",
//			args: args{
//				url:   "http://localhost:8080/gauge/",
//				name:  "test",
//				value: " 100",
//			},
//			want: "http://localhost:8080/gauge/test/ 100",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := getURL(tt.args.url, tt.args.name, tt.args.value); got != tt.want {
//				t.Errorf("getURL() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_update(t *testing.T) {
//	type args struct {
//		memory  *runtime.MemStats
//		gauge   storage.Gauge
//		counter storage.Counter
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  storage.Gauge
//		want1 storage.Counter
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := update(tt.args.memory, tt.args.gauge, tt.args.counter)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("update() got = %v, want %v", got, tt.want)
//			}
//			if !reflect.DeepEqual(got1, tt.want1) {
//				t.Errorf("update() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_uploadCounter(t *testing.T) {
//	type args struct {
//		client *http.Client
//		url    string
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			uploadCounter(tt.args.client, tt.args.url)
//		})
//	}
//}
//
//func Test_uploadGauge(t *testing.T) {
//	type args struct {
//		client *http.Client
//		url    string
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			uploadGauge(tt.args.client, tt.args.url)
//		})
//	}
//}
