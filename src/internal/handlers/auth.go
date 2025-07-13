package handlers

import (
	"fmt"
	"net/http"
)

func SignIn(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "SignIn Handler")
}

func SignUp(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "SignUp Handler")
}

func LogOut(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "LogOut Handler")
}
