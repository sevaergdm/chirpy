package main

import (
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sevaergdm/chirpy/internal/auth"
)

type RefreshToken struct {
	Token string `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt time.Time `json:"revoked_at"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}


	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No refresh token present", err)
		return
	}

	dbToken, err := cfg.dbQueries.GetRefreshToken(req.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token not found", err)
		return
	}

	userID, err := cfg.dbQueries.GetUserFromRefreshToken(req.Context(), dbToken.Token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has expired or been revoked", err)
		return
	}

	jwtToken, err := auth.MakeJWT(userID, os.Getenv("JWT_SECRET"))	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create new auth token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{Token: jwtToken})
}
