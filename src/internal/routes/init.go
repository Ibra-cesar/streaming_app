package routes

import (
	"net/http"

	"github.com/Ibra-cesar/video-streaming/src/internal/handlers"
)

func RegisterPubRoutes(mux *http.ServeMux, authHandlers *handlers.AuthConnServices, refreshHandler *handlers.RefreshHandlers){
	AuthRouter(mux, authHandlers)
	RefreshRouter(mux, refreshHandler)
}

func RegisterPrivRoutes(mux *http.ServeMux){
	HomeRouter(mux)
}
