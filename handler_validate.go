package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned := checkForProfaneWords(params.Body)
	respBody := returnVals{CleanedBody: cleaned}

	respondWithJSON(w, http.StatusOK, respBody)
}

func checkForProfaneWords(body string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	result := []string{}
	for _, word := range strings.Split(body, " ") {
		for _, profaneWord := range profaneWords {
			if strings.ToLower(word) == profaneWord {
				word = "****"
			}
		}
		result = append(result, word)
	}
	return strings.Join(result, " ")
}
