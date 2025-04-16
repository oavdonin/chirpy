package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/oavdonin/chirpy/internal/auth"
	"github.com/oavdonin/chirpy/internal/database"
)

func (cfg *apiConfig) authenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	const default_token_expiration int = 3600
	type parameters struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		Expiration int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}
	genuineHash, err := cfg.db.GetUserHash(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User not found", err)
		return
	}
	err = auth.CheckPasswordHash(genuineHash, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password incorrect", err)
		return
	}
	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	if params.Expiration == 0 || params.Expiration > default_token_expiration {
		params.Expiration = default_token_expiration
	}
	token, err := auth.MakeJWT(user.ID, string(cfg.jwtSigningKey))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error occurred while creating JWT token", err)
		return
	}
	timestampTokenCreated := time.Now().UTC()
	timestampTokenExpired := timestampTokenCreated.Add(time.Hour * 24 * 60)
	refreshToken := auth.MakeRefreshToken()
	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: timestampTokenCreated,
		UpdatedAt: timestampTokenCreated,
		UserID:    user.ID,
		ExpiresAt: &timestampTokenExpired,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't created the refresh token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}
