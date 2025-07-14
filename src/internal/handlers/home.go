package handlers

import (
	"fmt"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	//Handle least specfic routes e.g root "/"
	if r.URL.Path != "/"{
		http.NotFound(w,r)
		return
	}
	fmt.Fprintln(w, "SDOASJDIOASJIODAWID")
}

func AboutPage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Welcome to about page")
}
