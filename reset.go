package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "cannot delete reset outside dev environment")
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting all users - reset cancelled")
		return
	}

	cfg.fileserverHits.Store(0)
	respondWithJSON(w, http.StatusOK, "Reset complete!")
}
