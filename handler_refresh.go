package main

import (
	"net/http"
	"time"

	"github.com/kairos4213/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	rfTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), rfTokenStr)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
