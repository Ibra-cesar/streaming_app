package middleware

import (
	"log"
	"net/http"
	"time"
)

type ExtendWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *ExtendWriter) WriteHeader(statusCode int){
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Loggers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
			
		wrapper := &ExtendWriter{
			ResponseWriter: w,
			statusCode: http.StatusOK,
		}

	  next.ServeHTTP(wrapper, r)

		log.Println(wrapper.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}
