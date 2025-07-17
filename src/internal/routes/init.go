package routes

import (
	"net/http"

	"github.com/Ibra-cesar/video-streaming/src/internal/handlers"
)

func RegisterRoutes(mux *http.ServeMux, authHandlers *handlers.AuthConnServices){
	HomeRouter(mux)
	AuthRouter(mux, authHandlers)
}
