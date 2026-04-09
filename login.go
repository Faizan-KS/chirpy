package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Faizan-KS/chirpy/internal/auth"
	"github.com/Faizan-KS/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	reqBody := requestBody{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 1. Authenticate User
	hash, err := cfg.database.Gethash(r.Context(), reqBody.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if _, err := auth.CheckPasswordHash(reqBody.Password, hash); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please enter the right password"))
		return
	}

	// 2. Get User ID
	userId, err := cfg.database.GetId(r.Context(), reqBody.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. Generate JWT
	expires := time.Hour
	token, err := auth.MakeJWT(userId, cfg.jwtSecret, expires)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 4. Create Refresh Token
	refreshToken, err := auth.MakeRefreshToken()
	rToken, err := cfg.database.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    userId,
		ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
		RevokedAt: sql.NullTime{Valid: false},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(token)
	json.NewEncoder(w).Encode(rToken.Token)
}
