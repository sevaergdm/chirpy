package main

import (
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/sevaergdm/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token in header", err)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, os.Getenv("JWT_SECRET"))
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Unable to validate token", err)
		return
	}

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to get chrip", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not the author of the chirp", err)
		return
	}

	err = cfg.dbQueries.DeleteChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete chrip", err)
	}
	w.WriteHeader(http.StatusNoContent)
}
