package main

import (
	"encoding/json"
	"net/http"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Body string `json:"body"`
	}

	type Response struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	requestBody := RequestBody{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}
	if len(requestBody.Body) > 400 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, Response{
		Valid: true,
	})
}
