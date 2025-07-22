package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)
//JSON WRAPPER TO HANDLE STRUCT
func HttpError(w http.ResponseWriter, response any, statusCode int){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
//VALIDATION TO VALIDATE STRUCT
func Validate[T any](r *http.Request)(*T, error){
	var validate = validator.New()
	var input T

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		return nil, fmt.Errorf("Invalid JSON %w", err)
	}

	saniteze(&input)

	err = validate.Struct(input)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &input, nil
}
//SANITIZING STRINS INPUT
func saniteze(input any) {
	switch v := input.(type) {
		case *SignInInput:
			v.Email  = strings.TrimSpace(strings.ToLower(v.Email))
		case *SignUpInput:
			v.Username = strings.TrimSpace(v.Username)
			v.Email = strings.TrimSpace(strings.ToLower(v.Email))
	}
}
