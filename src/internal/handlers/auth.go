package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ibra-cesar/video-streaming/src/internal"
	"github.com/Ibra-cesar/video-streaming/src/internal/query_repo"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthConnServices struct {
	Queries                *query_repo.Queries
	JwtSecrets             []byte
	JwtRefreshTokenSecrets []byte
}

func AuthServices(queries *query_repo.Queries, jwt []byte, jwtRefreshToken []byte) *AuthConnServices {
	return &AuthConnServices{
		Queries:                queries,
		JwtSecrets:             jwt,
		JwtRefreshTokenSecrets: jwtRefreshToken,
	}
}

func (auth *AuthConnServices) SignUp(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	input, err := internal.Validate[internal.SignUpInput](r)
	if err != nil {
		internal.HttpError(w, internal.Response{
			Message: "Validation Error",
			Error:   err.Error(),
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	hash, err := HashPassword(input.Password)
	if err != nil {
		internal.HttpError(w, internal.Response{
			Message: "Failed to HashPassword",
			Error:   err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	query, err := auth.Queries.InsertUser(reqCtx, query_repo.InsertUserParams{
		ID:           uuid.New(),
		Name:         input.Username,
		Email:        input.Email,
		PasswordHash: hash,
	})
	if err != nil {
		internal.HttpError(w, internal.Response{Message: "Queries Error", Error: err.Error()}, http.StatusInternalServerError)
	}

	accesTok, err := internal.GenerateJWT(query.ID.String(), query.Email, auth.JwtSecrets)
	if err != nil {
		internal.HttpError(w, internal.Response{
			Message: "Failed to generate access token",
			Error:   err.Error(),
		}, http.StatusInternalServerError)
	}
	accesTokAge := int((30 * time.Minute).Seconds())
	SetCookie(w, "Access token", accesTok, true, false, accesTokAge, http.SameSiteLaxMode)

	refreshTok, err := internal.GenerateRefreshToken(query.ID.String(), auth.JwtRefreshTokenSecrets)
	if err != nil {
		internal.HttpError(w, internal.Response{
			Message: "Failed to generate acces refresh token",
			Error:   err.Error(),
		}, http.StatusInternalServerError)
	}
	refreshTokAge := int((7 * 24 * time.Hour).Seconds())
	ExpiresAt := pgtype.Timestamptz{
		Time:  time.Now().Add(7 * 24 * time.Hour),
		Valid: true,
	}
	_, err = auth.Queries.InsertNewUserToken(reqCtx, query_repo.InsertNewUserTokenParams{
		Token:     refreshTok,
		UserID:    query.ID,
		ExpiresAt: ExpiresAt,
	})
	if err != nil {
		internal.HttpError(w, internal.Response{
			Message: "Failed to query database",
			Error:   err.Error(),
		}, http.StatusInternalServerError)
	}
	SetCookie(w, "Refresh token", refreshTok, true, false, refreshTokAge, http.SameSiteLaxMode)

	internal.HttpError(w, internal.Response{
		Message: "User Created",
		Data:    query,
	}, http.StatusCreated)
}

func (auth *AuthConnServices) SignIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "this is sign up routes")
}
func (auth *AuthConnServices) LogOut(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "this is log out routes")
}
