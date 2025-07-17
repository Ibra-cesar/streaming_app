package routes

import (
	"net/http"

	"github.com/Ibra-cesar/video-streaming/src/internal/handlers"
)

func AuthRouter(mux *http.ServeMux, authHandler *handlers.AuthConnServices){
	mux.HandleFunc("POST /sign-in", authHandler.SignIn)
	mux.HandleFunc("POST /sign-up", authHandler.SignUp)
	mux.HandleFunc("POST /log-out", authHandler.LogOut)
}
