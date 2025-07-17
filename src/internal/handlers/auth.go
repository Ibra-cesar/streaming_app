package handlers

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthConnServices struct{
	conn *pgxpool.Pool
	jwtSecrets []byte
}

func AuthServices(pool *pgxpool.Pool, jwt []byte) *AuthConnServices{
	return &AuthConnServices{
		conn: pool,
		jwtSecrets: jwt,
	}
}

func (auth *AuthConnServices) SignIn(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "this is sign in routes")
}
func (auth *AuthConnServices) SignUp(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "this is sign up routes")
}
func (auth *AuthConnServices) LogOut(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "this is log out routes")
} 
