package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to hash Passsword, %w", err)
	}
	return  base64.URLEncoding.EncodeToString(hashBytes), nil
}

func ComparePassword(inputPass string, HashPass string) (bool) {
	decodedHash, err := base64.URLEncoding.DecodeString(HashPass)
	if err != nil {
		log.Printf("Failed to decoded password: %v", err) 
		return false
	}

	err = bcrypt.CompareHashAndPassword(decodedHash, []byte(inputPass))
	if err != nil {
		log.Printf("Password Didn't Match: %v", err)
		return  false
	}
	 log.Println("Password Match, Welcome")
	return true
}

func SetCookie(w http.ResponseWriter, name, token string, httpOnly bool, secure bool, age int, sameSite http.SameSite){
	cookie := &http.Cookie{
		Name: name,
		Value: token,
		HttpOnly: httpOnly,
		Secure: secure,
		SameSite: sameSite,
		MaxAge: age,
	}
	http.SetCookie(w, cookie)
}

func DeleteCookie(w http.ResponseWriter){
	cookie := &http.Cookie{
		Name: "Access-Token",
		Value: "",
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteDefaultMode,
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
