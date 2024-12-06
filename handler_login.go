package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kairos4213/chirpy/internal/auth"
	"github.com/kairos4213/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		AccessToken  string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding user info")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error finding user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshExpiresIn := time.Now().UTC().AddDate(0, 0, 60)
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = cfg.db.StoreRefreshToken(r.Context(), database.StoreRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: refreshExpiresIn,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
