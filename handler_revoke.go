package main

import (
	"net/http"

	"github.com/sevaergdm/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "No token included", err)
	}

	err = cfg.dbQueries.RevokeRefreshToken(req.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to revoke token", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
