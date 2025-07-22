package internal

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return  token.SignedString(secrets)
}

func ValidateToken[T jwt.Claims](token string, secrets []byte)(*T, error){
	var tokenStruct T
	t, err := jwt.ParseWithClaims(token, tokenStruct, func(t *jwt.Token) (any, error) {
		return secrets, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired){
			return nil, fmt.Errorf("%w",err)
		}
		if errors.Is(err, jwt.ErrTokenInvalidIssuer){
			return nil, fmt.Errorf("%w",err)
		}
		return nil, fmt.Errorf("Failed to parse Token, %w", err)
	}

	claims, ok := t.Claims.(T)
	if !ok || !t.Valid{
		return nil, fmt.Errorf("%w", jwt.ErrTokenInvalidClaims)
	}

	return &claims, nil
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
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return  token.SignedString(secrets)
}
