package main

import (
	"encoding/json"
	"http_adv/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}
	tokenCheck, err := cfg.database.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}
	if time.Now().After(tokenCheck.ExpiresAt) {
		http.Error(w, "refresh token expired", http.StatusUnauthorized)
		return
	}
	if tokenCheck.RevokedAt.Valid {
		http.Error(w, "refresh token revoked", http.StatusUnauthorized)
		return
	}
	expiresIn := time.Hour
	newAccessToken, err := auth.MakeJWT(tokenCheck.UserID, cfg.jwtSecret, expiresIn)
	if err != nil {
		http.Error(w, "failed to create token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newAccessToken)
}
