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

type RefreshHandlers struct{
	Q query_repo.Queries
	AccessToken []byte
	RefreshToken []byte
}

func RefreshService(q query_repo.Queries, accessTok []byte, refreshTok []byte) (*RefreshHandlers) {
	return &RefreshHandlers{
		Q: q,
		AccessToken: accessTok,
		RefreshToken: refreshTok,
	}
}

func(ref *RefreshHandlers) Refresh(w http.ResponseWriter, r *http.Request){

	//take context and validate cookie from request
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

	//validate the refresh token from cookie
	refreshTokenValid, err := internal.ValidateToken[internal.JwtRefreshTokenClaims](refreshToken, ref.RefreshToken)
	if err != nil{
		log.Printf("Invalid Token, %v", err)
		DeleteCookie(w, "refresh-token")
		internal.HttpError(w, internal.Response{
			Message: "Invalid Token",
			Error: err.Error(),
		}, http.StatusUnauthorized)
		return
	}

	//fetch refresh token from database
	storedRefToken, err := ref.Q.GetToken(reqCtx,refreshToken)
	if err != nil {
		if err == pgx.ErrNoRows{
			log.Printf("Refresh token not found: %v",err)
			internal.HttpError(w, internal.Response{
				Message: "Missing or Revoked Refresh Token",
				Error: err.Error(),
			}, http.StatusUnauthorized)
			return
		}
		log.Printf("Failed to get Token: %v", err)
		DeleteCookie(w, "refresh-token")
		internal.HttpError(w, internal.Response{
			Message: "Query Failed",
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	//check userId payload between client refresh token and server refresh token
	if parse, _ := uuid.Parse(refreshTokenValid.UserId); storedRefToken.UserID != parse {
		log.Printf("RefreshToken error: UserID mismatch: %s vs %s", refreshTokenValid.UserId, storedRefToken.UserID)
		DeleteCookie(w, "refresh-token")
		internal.HttpError(w, internal.Response{
			Message: "Invalid refresh token",
		}, http.StatusUnauthorized)
		return
	}

	//Revoked refresh token
	err = ref.Q.DeleteToken(reqCtx, storedRefToken.UserID) 
	if err != nil {
		log.Printf("Failed to delete token from db")
	}

	//fetch user information(id, Email)
	user, err := ref.Q.GetUser(reqCtx, storedRefToken.UserID)
	if err != nil {
		if err != pgx.ErrNoRows{
			log.Printf("User not found")
			DeleteCookie(w, "refresh-token")
			internal.HttpError(w, internal.Response{
				Message: "User is not found or invalid token",
				Error: pgx.ErrNoRows.Error(),
			}, http.StatusNotFound)
			return
		}
		log.Printf("Failed to fetch user")
		DeleteCookie(w, "refresh-token")
		internal.HttpError(w, internal.Response{
			Message: "Failed to fetch user",
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	//pass user info to generate new access token
	accessTok, err := internal.GenerateJWT(user.ID.String(), user.Email, ref.AccessToken) 
	if err != nil {
		internal.HttpError(w, internal.Response{
			Message: "Failed to generate token",
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	//pass user info to generate new refresh token
	refreshTok, err := internal.GenerateRefreshToken(user.ID.String(), ref.RefreshToken)
	if err != nil {
		internal.HttpError(w, internal.Response{
			Message: "Failed to generate token",
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	refreshTokAge := int((7 * 24 * time.Hour).Seconds())
	accessTokAge := int((30 * time.Minute).Seconds())
	tokenRow, err := ref.Q.InsertNewUserToken(reqCtx, query_repo.InsertNewUserTokenParams{
		Token: refreshTok,
		UserID: user.ID,
		ExpiresAt: pgtype.Timestamptz{
			Time: time.Now().Add(7 * 24 * time.Hour), 
			Valid: true,
		},
	})

	//Set new cookie for each token 
  SetCookie(w, "access-token", accessTok, true, false, accessTokAge, http.SameSiteLaxMode)
	SetCookie(w, "refresh-token", refreshTok, true, false, refreshTokAge, http.SameSiteLaxMode)

	//Response
	internal.HttpError(w, internal.Response{
		Message: "Token Refreshed",
		Data: map[string]string{
			"refresh-token": tokenRow.Token,
			"access-token": accessTok,
		},
	}, http.StatusCreated)
}
