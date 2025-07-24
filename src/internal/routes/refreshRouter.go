package routes

import (
	"net/http"

	"github.com/Ibra-cesar/video-streaming/src/internal/handlers"
)

func RefreshRouter(mux *http.ServeMux, refreshHandler *handlers.RefreshHandlers) {
	http.HandleFunc("POST /auth/refresh", refreshHandler.Refresh)
}
