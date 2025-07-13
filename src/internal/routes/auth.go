package routes

import (
	"net/http"

	"github.com/Ibra-cesar/video-streaming/src/internal/handlers"
)

func AuthRouter(mux *http.ServeMux){
	mux.HandleFunc("POST /sign-in", handlers.SignIn)
	mux.HandleFunc("POST /sign-up", handlers.SignUp)
	mux.HandleFunc("POST /log-out", handlers.LogOut)

	mux.HandleFunc("GET /sign-in", handlers.SignIn)
	mux.HandleFunc("GET /sign-up", handlers.SignUp)
	mux.HandleFunc("GET /log-out", handlers.LogOut)
}
