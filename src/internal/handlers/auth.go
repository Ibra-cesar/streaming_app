package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/Ibra-cesar/video-streaming/src/internal"
	"github.com/Ibra-cesar/video-streaming/src/internal/query_repo"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

//SignUp Handler
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
	SetCookie(w, "access-token", accesTok, true, false, accesTokAge, http.SameSiteLaxMode)

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
	SetCookie(w, "refresh-token", refreshTok, true, false, refreshTokAge, http.SameSiteLaxMode)

	internal.HttpError(w, internal.Response{
		Message: "User Created",
		Data:    query,
	}, http.StatusCreated)
}

//SignIn handler
func (auth *AuthConnServices) SignIn(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	input, err := internal.Validate[internal.SignInInput](r)
	if err != nil {
		log.Printf("Validation error: %v", err)
		internal.HttpError(w, internal.Response{
			Message: "Validation Error",
			Error: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	userInfo, err := auth.Queries.GetUserByEmail(reqCtx, input.Email)
	if err != nil {
		log.Printf("Email is not exist  exist")
		internal.HttpError(w, internal.Response{
			Message: "User is not exist exist please Sign Up",
			Error: pgx.ErrNoRows.Error(),
		}, http.StatusBadRequest)
		return
	}
	
	valid := ComparePassword(input.Password, userInfo.PasswordHash)
	if valid == false {
		log.Printf("Invalid Password")
		internal.HttpError(w, internal.Response{
			Message: "Wron Password",
		}, http.StatusBadRequest)
		return
	}

	DeleteCookie(w, "access-token")
	DeleteCookie(w, "refresh-token")

	err = auth.Queries.DeleteToken(reqCtx, userInfo.ID)
	if err != nil {
		log.Printf("Failed to revoked token: %v", err)
	}

	accessTokAge := int((30 * time.Minute).Seconds())
	refreshTokAge := int((7 * 24 * time.Hour).Seconds())
	ExpiresAt := pgtype.Timestamptz{
		Time:  time.Now().Add(7 * 24 * time.Hour),
		Valid: true,
	}

	accessToken, err := internal.GenerateJWT(userInfo.ID.String(), userInfo.Email, auth.JwtSecrets)
	if err != nil {
		log.Printf("Failed to Generate Acces Token")
		internal.HttpError(w, internal.Response{
			Message: "Failed to generate access token",
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	refreshToken, err := internal.GenerateRefreshToken(userInfo.ID.String(), auth.JwtRefreshTokenSecrets)
	if err != nil {
		log.Printf("Failed to Generate Refresh Token")
		internal.HttpError(w, internal.Response{
			Message: "Failed to Generate Refresh Token",
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	SetCookie(w, "access-token", accessToken, true, false, accessTokAge, http.SameSiteLaxMode)
	SetCookie(w, "refresh-token", refreshToken, true, false, refreshTokAge, http.SameSiteLaxMode)

	_, err = auth.Queries.InsertNewUserToken(reqCtx, query_repo.InsertNewUserTokenParams{
		Token: refreshToken,
		UserID: userInfo.ID,
		ExpiresAt: ExpiresAt,
	})
	if err != nil {
		log.Printf("Failed To insert Refresh Token")
		internal.HttpError(w, internal.Response{
			Message: "Failed to Insert refresh Token",
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	internal.HttpError(w, internal.Response{
		Message: "Login Successfull, Welcome!",
		Data: map[string]string{
			"name": userInfo.Name,
			"user-id": userInfo.ID.String(),
			"email": userInfo.Email,
		},
	}, http.StatusOK)
}

//LogOut Handler
func (auth *AuthConnServices) LogOut(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		if err == http.ErrNoCookie{
			internal.HttpError(w, internal.Response{
				Message: "Cookie not found, Unauthorized",
				Error: err.Error(),
			}, http.StatusUnauthorized)
			return
		}
		log.Printf("Failed to read cookie from request, %v", err)
		internal.HttpError(w, internal.Response{
			Message: "Error when read cookie",
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	 }
	 refreshToken := cookie.Value

	 tokenClaims, err := internal.ValidateToken[internal.JwtRefreshTokenClaims](refreshToken, auth.JwtRefreshTokenSecrets)
	 if err != nil{
		log.Printf("Invalid Token, %v", err)
		DeleteCookie(w, "refresh-token")
		internal.HttpError(w, internal.Response{
			Message: "Invalid Token",
			Error: err.Error(),
		}, http.StatusUnauthorized)
		return
	}

	userId, _ := uuid.Parse(tokenClaims.UserId)
	
	DeleteCookie(w, "access-token")
	DeleteCookie(w, "refresh-token")

	userInfo, err := auth.Queries.GetUser(reqCtx, userId)
	if err != nil {
		log.Printf("Log out failed failed to fetch user info")
		internal.HttpError(w, internal.Response{
			Message: "Log out failed",
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	_ = auth.Queries.DeleteToken(reqCtx, userInfo.ID)

	internal.HttpError(w, internal.Response{
		Message: "Log out successfull",
		Data: map[string]string{
			"user-id": userInfo.ID.String(),
			"name": userInfo.Name,
			"email": userInfo.Email,
		},
	}, http.StatusOK)
}
