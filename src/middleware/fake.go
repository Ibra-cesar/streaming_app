package middleware

import (
	"log"
	"net/http"
)

func FakeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		next.ServeHTTP(w, r)
		log.Println("Executing Fake Middleware")
	})
}
