package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/sevaergdm/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUserUpgrade(w http.ResponseWriter, req *http.Request) {
	const userUpgraded = "user.upgraded"

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to retrieve API key", err)
		return
	}

	if apiKey != os.Getenv("POLKA_KEY") {
		respondWithError(w, http.StatusUnauthorized, "Invalid API Key", err)
		return
	}

	var params parameters
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode payload", err)
		return
	}

	if params.Event != string(userUpgraded) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.dbQueries.UpgradeUser(req.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Unable to find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Unable to upgrade user", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
