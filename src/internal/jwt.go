package internal

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(userID, email string, secrets []byte) (string, error){
	claims := JwtClaims{
		UserID: userID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "http://localhost:8080/auth",
			Audience: []string{"http://localhost:8080/api/"},
			Subject: userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ID: uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return  token.SignedString(secrets)
}

func ValidateToken(token string, secrets []byte, claims jwt.Claims) (*jwt.Token, error) {
	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		return secrets, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("%w", err)
		}
		if errors.Is(err, jwt.ErrTokenInvalidIssuer) {
			return nil, fmt.Errorf("%w", err)
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !tkn.Valid {
		return nil, fmt.Errorf("%w", jwt.ErrTokenInvalidClaims)
	}

	return tkn, nil
}
func GenerateRefreshToken(userId string, secrets []byte)(string, error){
	claims := JwtClaims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "http://localhost:8080/auth",
			Audience: []string{"http://localhost:8080/api/"},
			Subject: userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ID: uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return  token.SignedString(secrets)
}
