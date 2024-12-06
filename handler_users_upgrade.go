package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kairos4213/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUsersUpgrade(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if key != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "API key is invalid")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var params parameters
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding user info")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "")
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = cfg.db.SetUserChirpyRedTrue(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	var resp struct{}
	respondWithJSON(w, http.StatusNoContent, resp)
}
