package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sevaergdm/chirpy/internal/auth"
	"github.com/sevaergdm/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUserUpdate(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token in request", err)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, os.Getenv("JWT_SECRET"))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User token not found", err)
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var params parameters
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to decode body", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to hash password", err)
		return
	}

	updatedUser, err := cfg.dbQueries.UpdateUserPassword(req.Context(), database.UpdateUserPasswordParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:           updatedUser.ID,
		CreatedAt:    updatedUser.CreatedAt,
		UpdatedAt:    updatedUser.UpdatedAt,
		Email:        updatedUser.Email,
		RefreshToken: bearerToken,
		IsChirpyRed:  updatedUser.IsChirpyRed,
	})
}
