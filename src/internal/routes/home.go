package routes

import (
	"net/http"

	"github.com/Ibra-cesar/video-streaming/src/internal/handlers"
)	

func HomeRouter(mux *http.ServeMux){
	mux.HandleFunc("GET /", handlers.HomePage)
	mux.HandleFunc("GET /about", handlers.AboutPage)
}
