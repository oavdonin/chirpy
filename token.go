package main

import (
	"net/http"
	"time"

	"github.com/oavdonin/chirpy/internal/auth"
	"github.com/oavdonin/chirpy/internal/database"
)

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is malformed", err)
		return
	}
	tokenAttrs, err := cfg.db.GetRefreshTokenAttrs(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error occured while validating a refresh token", err)
		return
	}

	now := time.Now().UTC()
	if tokenAttrs.RevokedAt != nil || tokenAttrs.ExpiresAt.Before(now) {
		respondWithError(w, http.StatusUnauthorized, "Token expired or revoked", err)
		return
	}

	token, err := auth.MakeJWT(tokenAttrs.UserID, string(cfg.jwtSigningKey))

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (cfg *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is malformed", err)
		return
	}
	tokenAttrs, err := cfg.db.GetRefreshTokenAttrs(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error occured while validating a refresh token", err)
		return
	}

	now := time.Now().UTC()
	if tokenAttrs.RevokedAt != nil || tokenAttrs.ExpiresAt.Before(now) {
		respondWithError(w, http.StatusUnauthorized, "Token expired or revoked", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		Token:     refreshToken,
		RevokedAt: &now,
		UpdatedAt: now,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error occured while revoking the refresh token", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)

	return
}
