package routes

import "net/http"

func RegisterRoutes(mux *http.ServeMux){
	HomeRouter(mux)
	AuthRouter(mux)
}
