package main

import (
	"net/http"

	"github.com/kairos4213/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	rfTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), rfTokenStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var resp interface{}
	respondWithJSON(w, http.StatusNoContent, resp)
}
