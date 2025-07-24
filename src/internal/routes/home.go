package routes

import (
	"net/http"

	"github.com/Ibra-cesar/video-streaming/src/internal/handlers"
	"github.com/Ibra-cesar/video-streaming/src/middleware"
)	

func HomeRouter(mux *http.ServeMux){
	mux.Handle("GET /", middleware.AuthMiddleware(http.HandlerFunc(handlers.HomePage)))
	mux.Handle("GET /about", middleware.AuthMiddleware(http.HandlerFunc(handlers.AboutPage)))
}
