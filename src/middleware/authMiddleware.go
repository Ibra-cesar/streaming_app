package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Ibra-cesar/video-streaming/src/internal"
	"github.com/joho/godotenv"
)

type ctxKey string

var UserCtxKey ctxKey = "user"

func env(name string) string {
	if name == "" {
		log.Fatal("Missing env variable name")
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to Load Env")
	}
	env := os.Getenv(name)
	return env
}


func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenSecret := []byte(env("JWT_SECRET_KEY"))

		header := r.Header.Get("Authorization")

		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			log.Printf("Missing Authorization Header")
			internal.HttpError(w, internal.Response{
				Message: "Missing Authorization Header",
			}, http.StatusUnauthorized)
			return 
		}

		accessToken := strings.TrimPrefix(header, "Bearer ")

		claims := &internal.JwtClaims{}
		_, err := internal.ValidateToken(accessToken, accessTokenSecret, claims)
		if err != nil {
			log.Printf("Invalid Token, %v, accessToken:%v", err, accessToken)
				internal.HttpError(w, internal.Response{
			Message: "Invalid Token",
			Error:   err.Error(),
		}, http.StatusUnauthorized)
		return
		}

		ctx := context.WithValue(r.Context(), UserCtxKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserClaims(ctx context.Context) (*internal.JwtClaims, error) {
	claims, ok := ctx.Value(UserCtxKey).(*internal.JwtClaims)
	if !ok {
		return nil, errors.New("Invalid claims")
	}
	return  claims, nil
}
