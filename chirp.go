package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oavdonin/chirpy/internal/database"
)

type RequestBody struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
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

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestBody := RequestBody{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}
	if (len(requestBody.Body) > 140 || len(requestBody.Body) == 0) && (len(requestBody.UserID) == 0) {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long or zero.", nil)
		return
	}
	replaceProfaneWords(&requestBody)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   requestBody.Body,
		UserID: requestBody.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error occured while chirp creation.", nil)
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
