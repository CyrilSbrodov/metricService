package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

//создание структуры gzip сжатия
type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {

	return w.Writer.Write(b)
}

//создание middleware gzip
func gzipHandle(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(rw, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(rw, gzip.BestSpeed)
		if err != nil {
			io.WriteString(rw, err.Error())
			return
		}
		defer gz.Close()

		rw.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: rw, Writer: gz}, r)
	}
}
