package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type RequestBody struct {
	Body string `json:"body"`
}

var badwords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func isBadWord(word string) bool {
	_, exists := badwords[word]
	return exists
}

func replaceProfaneWords(chirp *RequestBody) {
	s := strings.Split(chirp.Body, " ")
	for i, word := range s {
		if isBadWord(strings.ToLower(word)) {
			s[i] = "****"
		}
	}
	chirp.Body = strings.Join(s, " ")
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Valid       bool   `json:"valid"`
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	requestBody := RequestBody{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}
	if len(requestBody.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	replaceProfaneWords(&requestBody)
	respondWithJSON(w, http.StatusOK, Response{
		Valid:       true,
		CleanedBody: requestBody.Body,
	})
}
