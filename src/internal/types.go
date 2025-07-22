package internal

import (
	"github.com/golang-jwt/jwt/v5"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type SignInInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignUpInput struct {
	Username string `json:"username" validate:"required,min=5,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=10"`
}

type JwtClaims struct {
	UserID string
	Email  string
	jwt.RegisteredClaims
}

type JwtRefreshTokenClaims struct {
	UserId string
	jwt.RegisteredClaims
}
