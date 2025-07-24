package routes

import (
	"net/http"

	"github.com/Ibra-cesar/video-streaming/src/internal/handlers"
)

func AuthRouter(mux *http.ServeMux, authHandler *handlers.AuthConnServices){
	mux.HandleFunc("POST /auth/sign-in", authHandler.SignIn)
	mux.HandleFunc("POST /auth/sign-up", authHandler.SignUp)
	mux.HandleFunc("POST /auth/log-out", authHandler.LogOut)
}
